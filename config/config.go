package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	TelegramBotToken string
	AllowedUsersId []string

}

func LoadConfig() (*Config, error) {
    err := godotenv.Load(".env")
	if err != nil{
		log.Fatalf("Error loading .env file on main")
	}
	
	tgToken := os.Getenv("TELEGRAM_BOT_TOKEN")
	allowedIds := []string{
		os.Getenv("AllowedId1"),
		os.Getenv("AllowedId2"),
	}
		
	return &Config{
		TelegramBotToken: tgToken,
		AllowedUsersId: allowedIds,
	}, nil
}
