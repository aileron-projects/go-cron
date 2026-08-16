package cron

import (
	"context"
	"io"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/aileron-projects/go-tester"
)

func TestNewCron(t *testing.T) {
	t.Parallel()
	t.Run("cron parse error", func(t *testing.T) {
		cron, err := NewCron(&Config{
			Crontab: "INVALID",
			JobFunc: func(ctx context.Context) error { return nil },
		})
		tester.AssertEqual(t, nil, cron)
		tester.AssertEqualErr(t, &ParseError{What: "number of fields"}, err)
	})
}

func TestCron(t *testing.T) {
	t.Parallel()
	t.Run("run a job", func(t *testing.T) {
		var wg sync.WaitGroup
		wg.Add(1)
		var count atomic.Int32
		cron, _ := NewCron(&Config{
			Crontab: "* * * * * *",
			JobFunc: func(ctx context.Context) error {
				count.Add(1)
				wg.Done()
				return nil
			},
		})
		go cron.Start()
		wg.Wait()
		cron.Stop()
		tester.AssertEqual(t, 1, count.Load())
	})
	t.Run("already running", func(t *testing.T) {
		var wg sync.WaitGroup
		wg.Add(1)
		var count atomic.Int32
		cron, _ := NewCron(&Config{
			Crontab: "* * * * * *",
			JobFunc: func(ctx context.Context) error {
				defer wg.Done()
				count.Add(1)
				time.Sleep(time.Second)
				return nil
			},
		})
		go cron.Start()
		go cron.Start()
		wg.Wait()
		cron.Stop()
		tester.AssertEqual(t, 1, count.Load())
	})
	t.Run("calibrate", func(t *testing.T) {
		var wg sync.WaitGroup
		wg.Add(1)
		var count atomic.Int32
		cron, _ := NewCron(&Config{
			Crontab: "*/10 * * * *",
			JobFunc: func(ctx context.Context) error {
				count.Add(1)
				wg.Done()
				return nil
			},
		})
		now := time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC)
		cron.WithTimeFunc(func() time.Time {
			v := now
			now = now.Add(9*time.Minute + 59*time.Second)
			return v
		})
		durations := []time.Duration{}
		cron.WithTimeAfterFunc(func(d time.Duration) <-chan time.Time {
			durations = append(durations, d)
			return time.After(time.Second)
		})
		go cron.Start()
		wg.Wait()
		cron.Stop()
		tester.AssertEqual(t, 1, count.Load())
		tester.AssertEqual(t, 9*time.Minute+30*time.Second, durations[0])
		tester.AssertEqual(t, time.Second, durations[1])
	})
}

func TestRunner(t *testing.T) {
	t.Parallel()
	t.Run("queue=1", func(t *testing.T) {
		var wg sync.WaitGroup
		var count atomic.Int32
		r := &runner{
			queue: make(chan struct{}, 1),
			jobFunc: func(_ context.Context) error {
				defer wg.Done()
				count.Add(1)
				time.Sleep(time.Second)
				return nil
			},
			eventFunc: func(e Event, a ...any) {
				if e == EventJobAccepted {
					wg.Add(1)
				}
			},
		}
		r.Run() // Run 1. Should run
		r.Run() // Run 2. Should not run
		r.Run() // Run 3. Should not run
		wg.Wait()
		tester.AssertEqual(t, 1, count.Load())
	})
	t.Run("queue=2", func(t *testing.T) {
		var wg sync.WaitGroup
		var count atomic.Int32
		r := &runner{
			queue: make(chan struct{}, 2),
			jobFunc: func(_ context.Context) error {
				defer wg.Done()
				count.Add(1)
				time.Sleep(time.Second)
				return nil
			},
			eventFunc: func(e Event, a ...any) {
				if e == EventJobAccepted {
					wg.Add(1)
				}
			},
		}
		r.Run() // Run 1. Should run
		r.Run() // Run 2. Should run
		r.Run() // Run 3. Should not run
		wg.Wait()
		tester.AssertEqual(t, 2, count.Load())
	})
	t.Run("queue=3", func(t *testing.T) {
		var wg sync.WaitGroup
		var count atomic.Int32
		r := &runner{
			queue: make(chan struct{}, 3),
			jobFunc: func(_ context.Context) error {
				defer wg.Done()
				count.Add(1)
				time.Sleep(time.Second)
				return nil
			},
			eventFunc: func(e Event, a ...any) {
				if e == EventJobAccepted {
					wg.Add(1)
				}
			},
		}
		r.Run() // Run 1. Should run
		r.Run() // Run 2. Should run
		r.Run() // Run 3. Should run
		wg.Wait()
		tester.AssertEqual(t, 3, count.Load())
	})
	t.Run("job error", func(t *testing.T) {
		var wg sync.WaitGroup
		var err error
		r := &runner{
			queue: make(chan struct{}, 1),
			jobFunc: func(_ context.Context) error {
				return io.EOF // Job returns error.
			},
			eventFunc: func(e Event, a ...any) {
				switch e {
				case EventJobAccepted:
					wg.Add(1)
				case EventJobFailed:
					err = a[0].(error)
					wg.Done()
				}
			},
		}
		r.Run()
		wg.Wait()
		tester.AssertEqualErr(t, io.EOF, err)
	})
	t.Run("job panic", func(t *testing.T) {
		var wg sync.WaitGroup
		var err error
		r := &runner{
			queue: make(chan struct{}, 1),
			jobFunc: func(_ context.Context) error {
				panic(io.EOF) // Job panics.
			},
			eventFunc: func(e Event, a ...any) {
				switch e {
				case EventJobAccepted:
					wg.Add(1)
				case EventJobPanicked:
					err = a[0].(error)
					wg.Done()
				}
			},
		}
		r.Run()
		wg.Wait()
		tester.AssertEqualErr(t, io.EOF, err)
	})
	t.Run("run with context", func(t *testing.T) {
		var wg sync.WaitGroup
		r := &runner{
			queue: make(chan struct{}, 1),
			jobFunc: func(ctx context.Context) error {
				defer wg.Done()
				val := ctx.Value("foo")
				tester.AssertEqual(t, any("bar"), val)
				return nil
			},
			withContext: func() context.Context {
				return context.WithValue(context.Background(), "foo", "bar")
			},
		}
		r.Run()
		wg.Wait()
	})
}
