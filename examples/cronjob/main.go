package main

import (
	"context"
	"flag"
	"log"

	"github.com/aileron-projects/go-cron"
)

var (
	schedule = flag.String("schedule", "* * * * * *", "cron expression")
)

func main() {
	flag.Parse()
	log.Println("Crontab:", *schedule)

	count := 0
	c, err := cron.NewCron(&cron.Config{
		Crontab: *schedule,
		JobFunc: func(ctx context.Context) error {
			count++
			log.Println("Count:", count)
			return nil
		},
	})
	if err != nil {
		log.Panicln(err)
	}
	c.Start()
}
