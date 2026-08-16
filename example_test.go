package cron_test

import (
	"context"
	"fmt"
	"time"

	"github.com/aileron-projects/go-cron"
)

func ExampleCron() {
	now := time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC)
	count := 0
	stop := make(chan struct{})

	job, _ := cron.NewCron(&cron.Config{
		Crontab: "* * * * * *",
		JobFunc: func(ctx context.Context) error {
			if count++; count > 3 {
				stop <- struct{}{}
				return nil
			}
			now = now.Add(time.Second)
			fmt.Println(now.Format(time.DateTime), "Hello Go!")
			return nil
		},
	})

	job.WithTimeFunc(func() time.Time { return now })
	go job.Start()
	<-stop

	// Output:
	// 2000-01-01 00:00:01 Hello Go!
	// 2000-01-01 00:00:02 Hello Go!
	// 2000-01-01 00:00:03 Hello Go!
}

func ExampleCrontab_everySeconds() {
	ct, err := cron.Parse("* * * * * *") // Every seconds.
	if err != nil {
		panic(err)
	}

	// Replace internal clock for testing.
	now := time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC)
	ct.WithTimeFunc(func() time.Time { return now })

	for range 6 {
		fmt.Println(now.Format(time.DateTime), "|", ct.Next().Format(time.DateTime))
		now = now.Add(500 * time.Millisecond) // Forward 500ms.
	}
	// Output:
	// 2000-01-01 00:00:00 | 2000-01-01 00:00:01
	// 2000-01-01 00:00:00 | 2000-01-01 00:00:01
	// 2000-01-01 00:00:01 | 2000-01-01 00:00:02
	// 2000-01-01 00:00:01 | 2000-01-01 00:00:02
	// 2000-01-01 00:00:02 | 2000-01-01 00:00:03
	// 2000-01-01 00:00:02 | 2000-01-01 00:00:03
}

func ExampleCrontab() {
	ct, err := cron.Parse("59 59 23 31 * *") // Schedule at 31st at 23:59:59
	if err != nil {
		panic(err)
	}

	// Replace internal clock for testing.
	now := time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC)
	ct.WithTimeFunc(func() time.Time { return now })

	for range 10 {
		fmt.Println(now.Format(time.DateTime), "|", ct.Next().Format(time.DateTime))
		now = now.Add(15 * 24 * time.Hour) // Forward 15 days.
	}

	// Output:
	// 2000-01-01 00:00:00 | 2000-01-31 23:59:59
	// 2000-01-16 00:00:00 | 2000-01-31 23:59:59
	// 2000-01-31 00:00:00 | 2000-01-31 23:59:59
	// 2000-02-15 00:00:00 | 2000-03-31 23:59:59
	// 2000-03-01 00:00:00 | 2000-03-31 23:59:59
	// 2000-03-16 00:00:00 | 2000-03-31 23:59:59
	// 2000-03-31 00:00:00 | 2000-03-31 23:59:59
	// 2000-04-15 00:00:00 | 2000-05-31 23:59:59
	// 2000-04-30 00:00:00 | 2000-05-31 23:59:59
	// 2000-05-15 00:00:00 | 2000-05-31 23:59:59
}

func ExampleCrontab_NextAfter() {
	ct, err := cron.Parse("0 0 0 */5 * *") // Schedule every 5 days.
	if err != nil {
		panic(err)
	}

	// Replace internal clock for testing.
	now := time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)
	ct.WithTimeFunc(func() time.Time { return now })

	fmt.Println(ct.Next())
	fmt.Println(ct.NextAfter(time.Date(2026, 7, 31, 0, 0, 0, 0, time.UTC)))
	fmt.Println(ct.NextAfter(time.Date(2026, 8, 15, 0, 0, 0, 0, time.UTC)))
	fmt.Println(ct.NextAfter(time.Date(2026, 10, 01, 0, 0, 0, 0, time.UTC)))

	// Output:
	// 2026-01-06 00:00:00 +0000 UTC
	// 2026-08-01 00:00:00 +0000 UTC
	// 2026-08-16 00:00:00 +0000 UTC
	// 2026-10-06 00:00:00 +0000 UTC
}
