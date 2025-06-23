package main

import (
	"log"
	"net/http"

	"speedtest/internal/config"
	"speedtest/internal/server"
	"speedtest/internal/speedtester"
)

func main() {
	cfg := config.MustLoad()

	tester, err := speedtester.New(cfg.AppID, cfg.AppHash, cfg.BotToken)
	if err != nil {
		log.Fatalf("Не удалось создать тестер скорости: %s", err)
	}

	srv := &http.Server{
		Addr:    cfg.Port,
		Handler: server.Routes(tester),
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Не удалось запустить сервер: %s\n", err)
		}
	}()

	log.Printf("Сервер speedtest запущен на порту %s", cfg.Port)

	select {}
}
