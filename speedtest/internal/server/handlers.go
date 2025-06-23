package server

import (
	"encoding/json"
	"log"
	"net/http"
	"speedtest/internal/config"
	"speedtest/internal/speedtester"
	"strconv"
)

func handleSpeedTest(tester *speedtester.Tester, cfg *config.Config) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		query := r.URL.Query()
		chatIDStr := query.Get("chat_id")
		chatID, err := strconv.ParseInt(chatIDStr, 10, 64)
		if err != nil {
			http.Error(w, "invalid chat_id parameter", http.StatusBadRequest)
			return
		}

		appIDStr := query.Get("app_id")
		if appIDStr == "" {
			appIDStr = cfg.AppID
		}
		log.Println("appIDSTR: ", appIDStr)
		appID, err := strconv.Atoi(appIDStr)
		if err != nil {
			http.Error(w, "invalid app_id parameter", http.StatusBadRequest)
			return
		}

		appHash := query.Get("app_hash")
		if appHash == "" {
			appHash = cfg.AppHash
		}

		proxyAddress := query.Get("proxy")

		fileSizeMBStr := query.Get("mb")
		fileSizeMB, err := strconv.Atoi(fileSizeMBStr)
		if err != nil || fileSizeMB <= 0 {
			fileSizeMB = 10
		}

		log.Printf("Запуск теста скорости с размером файла %d MB для чата %d", fileSizeMB, chatID)

		result, err := tester.Measure(r.Context(), appID, appHash, proxyAddress, fileSizeMB, chatID)
		if err != nil {
			log.Printf("Ошибка при выполнении теста скорости: %s", err)
			http.Error(w, "internal server error", http.StatusInternalServerError)
			return
		}

		log.Printf("Тест скорости завершен. Загрузка: %.2f Мбит/с, Скачивание: %.2f Мбит/с", result.UploadSpeedMbps, result.DownloadSpeedMbps)

		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(result); err != nil {
			log.Printf("Ошибка при кодировании JSON-ответа: %s", err)
			http.Error(w, "internal server error", http.StatusInternalServerError)
		}
	}
}
