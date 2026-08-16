package main

import (
	"context"
	"errors"
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
			switch {
			case count%10 == 0: // Panic once in a 10 seconds.
				panic("panic in job")
			case count%5 == 0: // Error once in a 5 seconds.
				return errors.New("error from job")
			}
			return nil
		},
		EventHook: func(e cron.Event, a ...any) {
			log.Printf("  | %s: %+v", e, a)
		},
	})
	if err != nil {
		log.Panicln(err)
	}
	c.Start()
}
