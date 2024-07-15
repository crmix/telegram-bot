package item

import (
	"log"

	"github.com/robfig/cron/v3"
)


// Cron sets up a cron job to run the provided function every 30 seconds
func Cron(RunFunc func()) {
	c := cron.New(cron.WithSeconds())

	// Schedule the job to run every 30 seconds
	_, err := c.AddFunc("0 20 9 * * MON-FRI", RunFunc)
	if err != nil {
		log.Fatalf("Failed to add cron job: %v", err)
	}

	// Start the cron scheduler
	c.Start()

	// Keep the function running
	select {}
}
