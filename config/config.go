package config

import (
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	TelegramBotToken string
	OtabekAkaID      int64
	ElyorAkaID       int64
	FarruxAkaID      int64
	PsqlInfo         string
}

func LoadConfig() (*Config, error) {
	err := godotenv.Load(".env")
	if err != nil {
		log.Println("error loading env file on config package")
	}

	dbHost := os.Getenv("POSTGRES_HOST")
	dbPort := os.Getenv("POSTGRES_PORT")
	dbUser := os.Getenv("POSTGRES_USER")
	dbPassword := os.Getenv("POSTGRES_PASSWORD")
	dbName := os.Getenv("POSTGRES_DB")

	psqlInfo := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		dbHost, dbPort, dbUser, dbPassword, dbName)

	var tgToken string
	env := os.Getenv("ENVIRONMENT")
	if env == "dev" {
		tgToken = os.Getenv("TEST_TELEGRAM_BOT_TOKEN")
	} else if env == "prod" {
		tgToken = os.Getenv("TELEGRAM_BOT_TOKEN")
	}

	otabekakaid, err := strconv.ParseInt(os.Getenv("OtabekAkaId"), 10, 64)
	if err != nil {
		log.Printf("error during converting type in config file1")
	}

	elyorakaid, err := strconv.ParseInt(os.Getenv("ElyorAkaId"), 10, 64)
	if err != nil {
		log.Printf("error during converting type in config file2")
	}

	farruxakaid, err := strconv.ParseInt(os.Getenv("FarruxAkaId"), 10, 64)
	if err != nil {
		log.Printf("error during converting type in config file3")
	}

	return &Config{
		TelegramBotToken: tgToken,
		OtabekAkaID:      otabekakaid,
		ElyorAkaID:       elyorakaid,
		FarruxAkaID:      farruxakaid,
		PsqlInfo:         psqlInfo,
	}, nil
}
