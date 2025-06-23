package config

import (
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	TgToken   string
	TgHost    string
	SpeedTestHost string
	BatchSize int
}

func MustLoad() *Config {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, reading from environment variables")
	}

	token := os.Getenv("TELEGRAM_TOKEN")
	if token == "" {
		log.Fatal("TELEGRAM_BOT_TOKEN is not set")
	}

	host := os.Getenv("TELEGRAM_HOST")
	if host == "" {
		host = "api.telegram.org"
	}

	batchSizeStr := os.Getenv("BATCH_SIZE")
	if batchSizeStr == "" {
		batchSizeStr = "100"
	}

	speedTestHost := os.Getenv("SPEEDTEST_HOST")
	if speedTestHost == "" {
		speedTestHost = "localhost:8081"
	}

	batchSize, err := strconv.Atoi(batchSizeStr)
	if err != nil {
		log.Fatalf("invalid BATCH_SIZE: %v", err)
	}

	return &Config{
		TgToken:   token,
		TgHost:    host,
		BatchSize: batchSize,
		SpeedTestHost: speedTestHost,
	}
}
