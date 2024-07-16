package business

import (
	"fmt"
	"log"
	"telegram-bot/internal/database"

	"github.com/robfig/cron/v3"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type EmployeeService struct {
	repo *database.Repository
}

func NewEmployeeService(repo *database.Repository) *EmployeeService {
	return &EmployeeService{repo: repo}
}

func (s *EmployeeService) GetAndUpdateDutyEmployee() (database.Employee, error) {
	return s.repo.GetDutyEmployeeData()
}

func (s *EmployeeService) GetNextEmployee() (string, error) {
	return s.repo.GetNextDutyEmployee()
}

func (s *EmployeeService) RetrievingGroupID() (int64, error) {
	return s.repo.GettingGroupsId()
}

func (s *EmployeeService) InsertGroup(groupId int64) error {
	return s.repo.InsertGroupChatId(groupId)
}

func (s *EmployeeService) SendDailyDutyNotification(bot *tgbotapi.BotAPI) {

	c := cron.New()

	_, err := c.AddFunc("44 9 * * MON-FRI", func() {
		groupId, err := s.RetrievingGroupID()
		if err != nil {
			fmt.Printf("error during receiving groupid on api %v", err)
		}
		if err != nil {
			fmt.Printf("error during adding cron Addfunc %v", err)
		}

		employee, err := s.GetAndUpdateDutyEmployee()
		if err != nil {
			log.Printf("Failed to get and update duty employee: %v", err)
			return
		}

		remainderMsg := fmt.Sprintf("Xayrli tong! Bizning bugungi navbatchimiz: %s", employee.Name)
		msg := tgbotapi.NewMessage(groupId, remainderMsg)
		_, err = bot.Send(msg)
		if err != nil {
			log.Printf("Failed to send reminder msg: %s", err)
		}

	})
	if err != nil {
		log.Fatalf("Failed to add cron job: %v", err)
	}

	_, err = c.AddFunc("45 9 * * MON-FRI", func() {
		groupId, err := s.RetrievingGroupID()
		if err != nil {
			fmt.Printf("error during retrieving group ID: %v\n", err)
			return
		}
		employee, err := s.GetAndUpdateDutyEmployee()
		if err != nil {
			log.Printf("Failed to get and update duty employee: %v", err)
			return
		}
		resultsMsg := fmt.Sprintf("%s, Umid qilamanki bugungi navbatchiligingizni bajardingiz", employee.Name)
		msg := tgbotapi.NewMessage(groupId, resultsMsg)
		_, err = bot.Send(msg)
		if err != nil {
			log.Printf("Failed to send results message: %v", err)
		}
	})
	if err != nil {
		log.Fatalf("Failed to add cron job: %v", err)
	}

	c.Start()

	select {}
}
