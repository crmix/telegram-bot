package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	TelegramBotToken string
	AllowedUsersId []string
	AuthorizedUserID string

}

func LoadConfig() (*Config, error) {
    err := godotenv.Load(".env")
	if err != nil{
		log.Fatalf("Error loading .env file on main")
	}
	
	tgToken := os.Getenv("TEST_TELEGRAM_BOT_TOKEN")
	allowedIds := []string{
		os.Getenv("AllowedId1"),
		os.Getenv("AllowedId2"),
		os.Getenv("AllowedId3"),
	}
	authorisedId :=os.Getenv("AuthorisedUser")
		
	return &Config{
		TelegramBotToken: tgToken,
		AllowedUsersId: allowedIds,
		AuthorizedUserID: authorisedId,
	}, nil
}
