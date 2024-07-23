package api

import (
	"fmt"
	"log"
	"telegram-bot/config"
	"telegram-bot/internal/business"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type Bot struct {
	bot               *tgbotapi.BotAPI
	employeeService   *business.EmployeeService
	validationService *business.ValidationService
	GroupChatID       int64
}

var AllowedUsers map[int64]bool
var Id1, Id2, Id3 int64

func init() {
	var err error
	conf, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	Id1 = conf.OtabekAkaID
	Id2 = conf.ElyorAkaID
	Id3 = conf.FarruxAkaID

}

func StartBot(bot *tgbotapi.BotAPI, updates tgbotapi.UpdatesChannel, employeeService *business.EmployeeService, validationService *business.ValidationService, groupChatID int64) {

	bot.Debug = true
	b := &Bot{
		bot:               bot,
		employeeService:   employeeService,
		validationService: validationService,
		GroupChatID:       groupChatID,
	}

	for update := range updates {
		if update.MyChatMember != nil {
			b.handleMyChatMember(update.MyChatMember)
		}

		if update.Message != nil {
			userID := update.Message.From.ID
			if userID == Id1 || userID == Id2 || userID == Id3 && !update.Message.IsCommand() {
				forwardMessage(*b, bot, update.Message)
			}
			b.handleMessage(update.Message)
		}

		if update.CallbackQuery != nil {
			b.handleCallbackQuery(update.CallbackQuery)
		}

	}
}

func forwardMessage(b Bot, bot *tgbotapi.BotAPI, message *tgbotapi.Message) {
	groupChatID, err := b.employeeService.RetrievingGroupID()
	if err != nil {
		log.Printf("error during forwardMessage %v", err)
	}

	if message.Text != "" {
		msg := tgbotapi.NewMessage(groupChatID, message.Text)
		_, err := bot.Send(msg)
		if err != nil {
			log.Printf("Failed to forward text message: %v", err)
		}
	} else if message.Photo != nil {
		photo := message.Photo[len(message.Photo)-1]
		msg := tgbotapi.NewPhoto(groupChatID, tgbotapi.FileID(photo.FileID))
		_, err := bot.Send(msg)
		if err != nil {
			log.Printf("Failed to forward photo message: %v", err)
		}
	} else if message.Document != nil {
		doc := tgbotapi.NewDocument(groupChatID, tgbotapi.FileID(message.Document.FileID))
		_, err := bot.Send(doc)
		if err != nil {
			log.Printf("Failed to forward document message: %v", err)
		}
	} else if message.Video != nil {
		video := tgbotapi.NewVideo(groupChatID, tgbotapi.FileID(message.Video.FileID))
		_, err := bot.Send(video)
		if err != nil {
			log.Printf("Failed to forward video message: %v", err)
		}
	} else {
		log.Printf("Unhandled message type from authorized user: %v", message)
	}
}

func (b *Bot) handleMyChatMember(chatMember *tgbotapi.ChatMemberUpdated) {
	if chatMember.NewChatMember.User.UserName == b.bot.Self.UserName && chatMember.NewChatMember.Status == "member" {
		b.GroupChatID = chatMember.Chat.ID
		b.employeeService.InsertGroup(b.GroupChatID)
		welcomeMessage := "Assalomu aleykum! Men guruhda har kunlik navbatchilarni eslatib boradigan botman :)"
		msg := tgbotapi.NewMessage(b.GroupChatID, welcomeMessage)
		_, err := b.bot.Send(msg)
		if err != nil {
			log.Printf("Failed to send welcome msg: %s", err)
		}
	}
}

func (b *Bot) handleMessage(message *tgbotapi.Message) {

	if message.IsCommand() {
		switch message.Command() {
		case "start":
			b.handleStartCommand(message)
		case "validatephone":
			b.handleValidatePhoneCommand(message)
		case "tags":
			b.handleTagsCommand(message)
		case "restart":
			b.handleRestart()
		case "next":
			userId := message.From.ID
			if userId == Id1 || userId == Id2 || userId == Id3 {
				b.handleNextCommand()
			} else {
				b.bot.Send(tgbotapi.NewMessage(message.Chat.ID, "Sizga bu buyruqni amalga oshirishga ruxsat berilmagan"))
			}
		default:
			fmt.Println("Unknown command:", message.Command())
		}
	}
}

func (b *Bot) handleNextCommand() {
	groupid, err := b.employeeService.RetrievingGroupID()
	if err != nil {
		log.Printf("failed to get groupid for next tag %v", err)
	}
	nextEmployee, err := b.employeeService.GetNextEmployee()
	if err != nil {
		log.Printf("failed to get nextEmployee for next tag %v", err)
	}
	chat := fmt.Sprintf("Demak bugun navbatchilikni %s qiladilar", nextEmployee)
	msg := tgbotapi.NewMessage(groupid, chat)
	_, err = b.bot.Send(msg)
	if err != nil {
		log.Printf("Failed to send next msg %v", err)
	}
}

func (b *Bot) handleStartCommand(message *tgbotapi.Message) {
	msg := tgbotapi.NewMessage(message.Chat.ID, "Xush kelibsiz! /validatephone kamandasidan foydalanib raqamingizni validatsiya qiling.")
	_, err := b.bot.Send(msg)
	if err != nil {
		log.Printf("Failed to send start msg %v", err)
	}
}

func (b *Bot) handleValidatePhoneCommand(message *tgbotapi.Message) {
	phonenumber := message.CommandArguments()
	isValid := b.validationService.IsValidPhoneNumber(phonenumber)
	var responseMsg string
	if isValid {
		responseMsg = "Valid phone number."
	} else {
		responseMsg = "Invalid phone number."
	}
	msg := tgbotapi.NewMessage(message.Chat.ID, responseMsg)
	_, err := b.bot.Send(msg)
	if err != nil {
		log.Printf("Failed to send validation msg %s", err)
	}
}

func (b *Bot) handleRestart() {
	chatID, err := b.employeeService.RetrievingGroupID()
	if err != nil {
		log.Printf("Failed to get groupId on restart tag %v", err)
	}
	introductionMsg := "Salom, men qayta ishga tushdim! Men navbatchilikni eslatib turuvchi botman. "
	msg := tgbotapi.NewMessage(chatID, introductionMsg)
	_, err = b.bot.Send(msg)
	if err != nil {
		log.Printf("Failed to send restart introduction message: %v", err)
	}

	employee, err := b.employeeService.GetAndUpdateDutyEmployee()
	if err != nil {
		log.Printf("Failed to get and update duty employee on restart tag: %v", err)
		return
	}

	reminderMsg := fmt.Sprintf("Bugungi navbatchimiz: %s edi.", employee.Name)
	msg = tgbotapi.NewMessage(chatID, reminderMsg)
	_, err = b.bot.Send(msg)
	if err != nil {
		log.Printf("Failed to send restart remainder message: %v", err)
	}
}

func (b *Bot) handleTagsCommand(message *tgbotapi.Message) {
	tags := []string{"start", "validatephone", "restart"}
	inlineKeyboard := b.createTagsInlineKeyboard(tags)
	msg := tgbotapi.NewMessage(message.Chat.ID, "Please choose a tag:")
	msg.ReplyMarkup = inlineKeyboard
	b.bot.Send(msg)
}

func (b *Bot) createTagsInlineKeyboard(tags []string) tgbotapi.InlineKeyboardMarkup {
	var rows [][]tgbotapi.InlineKeyboardButton

	for _, tag := range tags {
		row := tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(tag, tag),
		)
		rows = append(rows, row)
	}

	return tgbotapi.NewInlineKeyboardMarkup(rows...)
}

func (b *Bot) handleCallbackQuery(callbackQuery *tgbotapi.CallbackQuery) {
	command := callbackQuery.Data

	switch command {
	case "start":
		b.handleStartCommand(callbackQuery.Message)
	case "validatephone":
		b.handleValidatePhoneCommand(callbackQuery.Message)
	case "restart":
		b.handleRestart()
	default:
		response := fmt.Sprintf("Unknown command: %s", command)
		msg := tgbotapi.NewMessage(callbackQuery.Message.Chat.ID, response)
		b.bot.Send(msg)
	}

	// Acknowledge the callback query
	callback := tgbotapi.NewCallback(callbackQuery.ID, "")
	if _, err := b.bot.Request(callback); err != nil {
		log.Printf("Failed to acknowledge callback query: %v", err)
	}
}
