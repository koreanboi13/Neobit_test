package telegram

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"strings"
)

const (
	TestCmd  = "/test"
	HelpCmd  = "/help"
	StartCmd = "/start"
)

func (p *Processor) doCmd(ctx context.Context, text string, chatID int) error {
	text = strings.TrimSpace(text)

	log.Printf("got new command '%s' from chat %d", text, chatID)

	parts := strings.Fields(text)
	command := parts[0]

	switch command {
	case TestCmd:
		go p.testSpeed(context.Background(), chatID, parts)
		return nil
	case HelpCmd:
		return p.sendHelp(ctx, chatID)
	case StartCmd:
		return p.sendHello(ctx, chatID)
	default:
		return p.tg.SendMessage(ctx, chatID, msgUnknownCommand)
	}
}

func (p *Processor) testSpeed(ctx context.Context, chatID int, args []string) error {
	if len(args) == 1 {
		if err := p.tg.SendMessage(ctx, chatID, msgSpeedTestStart); err != nil {
			return fmt.Errorf("failed to send starting message: %w", err)
		}

		result, err := p.speedtestClient.SpeedTest(int64(chatID), 0, "", "")
		if err != nil {
			log.Printf("speed test failed: %v", err)
			return p.tg.SendMessage(ctx, chatID, msgSpeedTestFailed)
		}

		responseText := fmt.Sprintf(
			"Speed Test Results:\n\nUpload: %.2f Мбит/с\nDownload: %.2f Мбит/с",
			result.UploadSpeedMbps,
			result.DownloadSpeedMbps,
		)

		return p.tg.SendMessage(ctx, chatID, responseText)

	} else if len(args) >= 3 {
		appID, err := strconv.Atoi(args[1])
		if err != nil {
			return p.tg.SendMessage(ctx, chatID, "Invalid App ID. It must be a number.")
		}

		appHash := args[2]

		var proxy string
		if len(args) > 3 {
			proxy = args[3]
		}
		if err := p.tg.SendMessage(ctx, chatID, msgSpeedTestStart); err != nil {
			return fmt.Errorf("failed to send starting message: %w", err)
		}

		result, err := p.speedtestClient.SpeedTest(int64(chatID), appID, appHash, proxy)
		if err != nil {
			log.Printf("speed test failed: %v", err)
			return p.tg.SendMessage(ctx, chatID, msgSpeedTestFailed)
		}

		responseText := fmt.Sprintf(
			"Speed Test Results:\n\nUpload: %.2f Мбит/с\nDownload: %.2f Мбит/с",
			result.UploadSpeedMbps,
			result.DownloadSpeedMbps,
		)

		return p.tg.SendMessage(ctx, chatID, responseText)
	} else {
		return p.tg.SendMessage(ctx, chatID, msgInvalidTestCmd)
	}
}

func (p *Processor) sendHelp(ctx context.Context, chatID int) error {
	return p.tg.SendMessage(ctx, chatID, msgHelp)
}

func (p *Processor) sendHello(ctx context.Context, chatID int) error {
	return p.tg.SendMessage(ctx, chatID, msgHello)
}
