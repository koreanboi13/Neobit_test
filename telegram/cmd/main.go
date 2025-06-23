package main

import (
	"context"
	"log"
	"time"

	"api/internal/clients/speedtest"
	tgClient "api/internal/clients/telegram"
	"api/internal/config"
	"api/internal/events/telegram"
)

func main() {
	cfg := config.MustLoad()

	telegramClient := tgClient.New(cfg.TgHost, cfg.TgToken)
	speedtestClient := speedtest.New(cfg.SpeedTestHost)

	eventProcessor := telegram.New(telegramClient, speedtestClient)

	log.Print("service started")

	fetcher := eventProcessor
	processor := eventProcessor

	go func() {
		for {
			events, err := fetcher.Fetch(context.Background(), cfg.BatchSize)
			if err != nil {
				log.Printf("error fetching events: %v", err)
				time.Sleep(1 * time.Second)
				continue
			}

			if len(events) == 0 {
				time.Sleep(1 * time.Second)
				continue
			}

			for _, event := range events {
				if err := processor.Process(context.Background(), event); err != nil {
					log.Printf("error processing event: %v", err)
				}
			}
		}
	}()

	select {}
}
