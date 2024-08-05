package main

import (
	"fmt"
	"log"
	"net/http"
	"telegram-bot/config"
	"telegram-bot/internal/api"
	"telegram-bot/internal/business"
	db "telegram-bot/internal/database"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func main() {

	tgToken, err := config.LoadConfig()

	if err != nil {
		log.Printf("err during recieving config file on main %v", err)
	}
	if tgToken.TelegramBotToken == "" {
		log.Fatalf("Telegram bot token is not set in config")
	}

	bot, err := tgbotapi.NewBotAPI(tgToken.TelegramBotToken)

	if err != nil {
		log.Println(err)
	}

	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	url := "https://api.telegram.org/bot" + tgToken.TelegramBotToken + "/getMe"
	resp, err := client.Get(url)
	if err != nil {
		log.Fatalf("HTTP request failed: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Fatalf("Non-OK HTTP status: %s", resp.Status)
	}

	log.Printf("Bot is authorized on account %s", bot.Self.UserName)

	updateConfig := tgbotapi.NewUpdate(0)
	updateConfig.Timeout = 60

	updates := bot.GetUpdatesChan(updateConfig)
	dbConn, err := db.NewDbConn()
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	repo := db.NewRepository(dbConn)
	employeeService := business.NewEmployeeService(repo)
	validationService := business.NewValidationService()

	fmt.Println("Bot is starting...")

	go api.StartBot(bot, updates, employeeService, validationService, 0)

	employeeService.SendDailyDutyNotification(bot)

}
