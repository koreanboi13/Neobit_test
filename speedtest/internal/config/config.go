package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	AppID    string
	AppHash  string
	BotToken string
	Port     string
}

func MustLoad() *Config {
	if err := godotenv.Load(); err != nil {
		log.Println(".env файл не найден")
	}

	appID := os.Getenv("APP_ID")
	if appID == "" {
		log.Fatal("APP_ID не установлен")
	}

	appHash := os.Getenv("APP_HASH")
	if appHash == "" {
		log.Fatal("APP_HASH не установлен")
	}

	botToken := os.Getenv("BOT_TOKEN")
	if botToken == "" {
		log.Fatal("BOT_TOKEN не установлен")
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = ":8081"
	} else if port[0] != ':' {
		port = ":" + port
	}

	return &Config{
		AppID:    appID,
		AppHash:  appHash,
		BotToken: botToken,
		Port:     port,
	}
}
