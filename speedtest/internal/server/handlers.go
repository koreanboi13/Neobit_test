package server

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"speedtest/internal/speedtester"
)

func handleSpeedTest(tester *speedtester.Tester) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		chatIDStr := r.URL.Query().Get("chat_id")
		chatID, err := strconv.ParseInt(chatIDStr, 10, 64)
		if err != nil {
			http.Error(w, "invalid chat_id parameter", http.StatusBadRequest)
			return
		}

		fileSizeMBStr := r.URL.Query().Get("mb")
		fileSizeMB, err := strconv.Atoi(fileSizeMBStr)
		if err != nil || fileSizeMB <= 0 {
			fileSizeMB = 10
		}

		log.Printf("Запуск теста скорости с размером файла %d MB для чата %d", fileSizeMB, chatID)

		result, err := tester.Measure(r.Context(), fileSizeMB, chatID)
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