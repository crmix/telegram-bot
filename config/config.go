package config

import (
	// "fmt"
	// "log"
	// "os"
	// "strconv"

	// "github.com/joho/godotenv"
)

type Config struct {
	TelegramBotToken string
	OtabekAkaID      int64
	ElyorAkaID       int64
	FarruxAkaID      int64
}

func LoadConfig() (*Config, error) {
	// err := godotenv.Load(".env")
	// if err != nil {
	// 	log.Println("error loading env file on config package")
	// }

	// dbHost := os.Getenv("DB_HOST")
	// dbPort := os.Getenv("DB_PORT")
	// dbUser := os.Getenv("DB_USER")
	// dbPassword := os.Getenv("DB_PASSWORD")
	// dbName := os.Getenv("DB_NAME")

	// psqlInfo := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
	// 	dbHost, dbPort, dbUser, dbPassword, dbName)

	// var tgToken string
	// if os.Getenv("ENVIRONMENT") == "dev" {
	// 	tgToken = os.Getenv("TEST_TELEGRAM_BOT_TOKEN")
	// } else if os.Getenv("ENVIROMENT") == "prod" {
	// 	tgToken = os.Getenv("TELEGRAM_BOT_TOKEN")
	// }

	// otabekakaid, err := strconv.ParseInt(os.Getenv("OtabekAkaId"), 10, 64)
	// if err != nil {
	// 	log.Printf("error during converting type in config file1")
	// }

	// elyorakaid, err := strconv.ParseInt(os.Getenv("ElyorAkaId"), 10, 64)
	// if err != nil {
	// 	log.Printf("error during converting type in config file2")
	// }

	// farruxakaid, err := strconv.ParseInt(os.Getenv("FarruxAkaId"), 10, 64)
	// if err != nil {
	// 	log.Printf("error during converting type in config file3")
	// }

	return &Config{
		TelegramBotToken: "7181421703:AAG3FmPqUgitSvbg8dwlUdGT8mbiMCRM1R8",
		OtabekAkaID:      150225351,
		ElyorAkaID:       262843262,
		FarruxAkaID:      1399785585,
	}, nil
}
