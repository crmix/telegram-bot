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

func (s *EmployeeService) GetPrevEmployee() (string, error) {
	return s.repo.GetPreviousDutyEmployee()
}

func (s *EmployeeService) RetrievingGroupID() (int64, error) {
	return s.repo.GettingGroupsId()
}

func (s *EmployeeService) InsertGroup(groupId int64) error {
	return s.repo.InsertGroupChatId(groupId)
}

func (s *EmployeeService) SendDailyDutyNotification(bot *tgbotapi.BotAPI) {

	c := cron.New()

	_, err := c.AddFunc("5 4 * * MON-FRI", func() {
		// fmt.Println("bot: ", bot)
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
		remainderMsg := fmt.Sprintf(`<b>Xayrli tong! Bizning bugungi navbatchimiz: <i>%s</i></b> 
		
		<b>Navbatchining vazifalari:</b>
		
		<b>Idishlarni yig‘ishtirish va tozalash</b> — Ishxona xonasida qolgan idishlarni yig‘ib, yuvilishi kerak bo‘lsa, yuvib qo‘yadi.
		<b>Zarur joylarni tozalash</b> — Agarda kimdir sizda qasti bo'lib tuzoq qo'ygan bo'lsa topib o'ziga aytib tozalattirish yoki o'zi tozalash.
		<b>Chiqindilarni tashlash</b> — Kun davomida yig‘ilgan chiqindilarni belgilangan joyga olib chiqib, axlat qutisini yangilaydi.
		<b>Texnikalarni o‘chirish</b> — Ish vaqti tugagach, kompyuter, printer va boshqa texnikalarni o‘chirish va xavfsizligini ta’minlash.
	  `, employee.Name)
		msg := tgbotapi.NewMessage(groupId, remainderMsg)
		msg.ParseMode = "HTML"
		_, err = bot.Send(msg)
		if err != nil {
			log.Printf("Failed to send reminder msg: %s", err)
		}

	})
	if err != nil {
		log.Fatalf("Failed to add cron job: %v", err)
	}

	_, err = c.AddFunc("20 9 * * MON-FRI", func() {
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
		resultsMsg := fmt.Sprintf("<b>%s</b>, Umid qilamanki bugungi navbatchiligingizni bajardingiz", employee.Name)
		msg := tgbotapi.NewMessage(groupId, resultsMsg)
		msg.ParseMode = "HTML"
		_, err = bot.Send(msg)
		if err != nil {
			log.Printf("Failed to send results message: %v", err)
		}
	})
	if err != nil {
		log.Fatalf("Failed to add cron job: %v", err)
	}

	_, err = c.AddFunc("0 11 * * MON-FRI", func() {
		groupId, err := s.RetrievingGroupID()
		if err != nil {
			fmt.Printf("error during retrieving group ID: %v\n", err)
			return
		}
		resultsMsg := fmt.Sprintf(`<b>Hurmatli hamkasblar!</b>
		<i>Oshxonamizni yillik muzeyga aylantirmaslik uchun, qolgan ovqatlar va idishlaringizni olib keting. Hidlari bilan san’at asariga aylantirish shart emas!</i>`)
		msg := tgbotapi.NewMessage(groupId, resultsMsg)
		msg.ParseMode = "HTML"
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
