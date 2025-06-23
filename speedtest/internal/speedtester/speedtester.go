package speedtester

import (
	"bytes"
	"context"
	"crypto/rand"
	"encoding/binary"
	"fmt"
	"io"
	"net/url"
	"time"

	"github.com/go-faster/errors"
	"github.com/gotd/td/telegram"
	"github.com/gotd/td/telegram/dcs"
	"github.com/gotd/td/telegram/downloader"
	"github.com/gotd/td/telegram/uploader"
	"golang.org/x/net/proxy"

	"github.com/gotd/td/tg"
)

type Tester struct {
	botToken string
}

func New(botToken string) (*Tester, error) {
	return &Tester{
		botToken: botToken,
	}, nil
}

func (t *Tester) resolvePeer(chatID int64) tg.InputPeerClass {
	return &tg.InputPeerUser{
		UserID: chatID,
	}
}

func (t *Tester) Measure(ctx context.Context, appID int, appHash string, proxyAddress string, fileSizeMB int, chatID int64) (*SpeedTestResult, error) {
	opts := telegram.Options{}
	if proxyAddress != "" {
		proxyURL, err := url.Parse(proxyAddress)
		if err != nil {
			return nil, errors.Wrap(err, "ошибка парсинга адреса прокси")
		}

		var auth *proxy.Auth
		if proxyURL.User != nil {
			password, _ := proxyURL.User.Password()
			auth = &proxy.Auth{
				User:     proxyURL.User.Username(),
				Password: password,
			}
		}

		dialer, err := proxy.SOCKS5("tcp", proxyURL.Host, auth, proxy.Direct)
		if err != nil {
			return nil, errors.Wrap(err, "ошибка создания SOCKS5 dialer")
		}

		dcDialer, ok := dialer.(proxy.ContextDialer)
		if !ok {
			return nil, errors.New("dialer не является proxy.ContextDialer")
		}

		opts.Resolver = dcs.Plain(dcs.PlainOptions{
			Dial: dcDialer.DialContext,
		})
	}

	//создание экземпляра клиента MTProto, для установки соединения с сервером Telegram
	//appID и appHash - уникальные значения, которые показывают телеграм, что клиент авторизирован
	//для их получения, необходимо было авторизоваться на my.telegram.org и зарегистрировать приложение
	client := telegram.NewClient(appID, appHash, opts)

	var result *SpeedTestResult
	var measureErr error

	//установка постоянного и зашифрованного соединения клиента с сервером телеграмм
	//соединение устанавливается только внутри функции
	if err := client.Run(ctx, func(ctx context.Context) error {
		//аутентификация клиента, который общается по MTProto c сервером телеграмм, как бота
		if _, err := client.Auth().Bot(ctx, t.botToken); err != nil {
			return errors.Wrap(err, "ошибка аутентификации бота")
		}

		uploadSpeed, fileLocation, fileSize, err := t.measureUpload(ctx, client, fileSizeMB, chatID)
		if err != nil {
			measureErr = errors.Wrap(err, "ошибка при измерении скорости загрузки")
			return measureErr
		}

		downloadSpeed, err := t.measureDownload(ctx, client, fileLocation, fileSize)
		if err != nil {
			measureErr = errors.Wrap(err, "ошибка при измерении скорости скачивания")
			return measureErr
		}

		result = &SpeedTestResult{
			UploadSpeedMbps:   uploadSpeed,
			DownloadSpeedMbps: downloadSpeed,
		}

		return nil
	}); err != nil {
		return nil, errors.Wrap(err, "ошибка выполнения клиента")
	}

	return result, measureErr
}

func (t *Tester) measureUpload(ctx context.Context, client *telegram.Client, fileSizeMB int, chatID int64) (float64, *tg.InputDocumentFileLocation, int64, error) {
	fileSize := fileSizeMB * 1024 * 1024
	data := make([]byte, fileSize)
	if _, err := rand.Read(data); err != nil {
		return 0, nil, 0, errors.Wrap(err, "не удалось сгенерировать случайные данные")
	}

	//создание загрузчика из gotd, который умеет загружать большие файлы путем разбиения их на части
	u := uploader.NewUploader(client.API())
	fileName := fmt.Sprintf("speedtest-%d.bin", time.Now().Unix())

	uploadStart := time.Now()
	upload, err := u.Upload(ctx, uploader.NewUpload(fileName, bytes.NewReader(data), int64(len(data))))
	if err != nil {
		return 0, nil, 0, errors.Wrap(err, "ошибка загрузки файла")
	}
	uploadDuration := time.Since(uploadStart)

	//рассчет скорости в Мбит/с: (размер_в_байтах * 8) / время_в_секундах / (1024*1024)
	uploadSpeed := (float64(len(data)) * 8 / uploadDuration.Seconds()) / (1024 * 1024)

	var randomID int64
	if err := binary.Read(rand.Reader, binary.LittleEndian, &randomID); err != nil {
		return 0, nil, 0, errors.Wrap(err, "не удалось сгенерировать случайный ID")
	}

	//загруженный файл отправляется пользователю, чтобы получить постоянную
	//ссылку на него и далее его можно было скачать
	updates, err := client.API().MessagesSendMedia(ctx, &tg.MessagesSendMediaRequest{
		Peer: t.resolvePeer(chatID),
		Media: &tg.InputMediaUploadedDocument{
			File:     upload,
			MimeType: "application/octet-stream",
			Attributes: []tg.DocumentAttributeClass{
				&tg.DocumentAttributeFilename{FileName: fileName},
			},
		},
		Message:  "Test File",
		RandomID: randomID,
	})
	if err != nil {
		return 0, nil, 0, errors.Wrap(err, "не удалось отправить медиа")
	}

	//ищем сообщение с нашим файлом, чтобы извлечь из него информацию для скачивания
	var document *tg.Document
	switch u := updates.(type) {
	case *tg.Updates:
		for _, update := range u.Updates {
			if msg, ok := update.(*tg.UpdateNewMessage); ok {
				if m, ok := msg.Message.(*tg.Message); ok {
					if media, ok := m.Media.(*tg.MessageMediaDocument); ok {
						if doc, ok := media.Document.(*tg.Document); ok {
							document = doc
							break
						}
					}
				}
			}
		}
	}

	if document == nil {
		return 0, nil, 0, errors.New("не удалось найти документ в ответе от сервера")
	}

	location := &tg.InputDocumentFileLocation{
		ID:            document.ID,
		AccessHash:    document.AccessHash,
		FileReference: document.FileReference,
	}

	return uploadSpeed, location, document.Size, nil
}

func (t *Tester) measureDownload(ctx context.Context, client *telegram.Client, location *tg.InputDocumentFileLocation, fileSize int64) (float64, error) {
	d := downloader.NewDownloader()

	downloadStart := time.Now()

	//скачивание файла, без сохранения файла
	if _, err := d.Download(client.API(), location).Stream(ctx, io.Discard); err != nil {
		return 0, errors.Wrap(err, "ошибка скачивания файла")
	}
	downloadDuration := time.Since(downloadStart)

	downloadSpeed := (float64(fileSize) * 8 / downloadDuration.Seconds()) / (1024 * 1024)

	return downloadSpeed, nil
}
