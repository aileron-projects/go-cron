<!-- markdownlint-disable MD033 MD041 -->

<div align="center">

[![Release](https://img.shields.io/github/v/release/aileron-projects/go-cron?sort=semver)](https://github.com/aileron-projects/go-cron/releases)
[![Reference](https://pkg.go.dev/badge/github.com/aileron-projects/go-cron.svg)](https://pkg.go.dev/github.com/aileron-projects/go-cron)
[![DeepWiki](https://deepwiki.com/badge.svg)](https://deepwiki.com/aileron-projects/go-cron)
[![Test](https://github.com/aileron-projects/go-cron/actions/workflows/test.yaml/badge.svg)](https://github.com/aileron-projects/go/actions/workflows/test.yaml)

[![Insights](https://badgen.net/badge/Insights/open%2Fsource%2Finsights/cyan)](https://deps.dev/go/github.com%2Faileron-projects%2Fgo-cron)
[![Insights](https://badgen.net/badge/Insights/OSS%2FInsight/orange)](https://ossinsight.io/analyze/aileron-projects/go-cron)

</div>

# go-cron

**Cron jobs scheduler and executer for Go.**

## Features

- Cron expression parser
- Rich cron expression
- Job schedule and run
- Zero dependency
- Lightweight

## Usages

### Crontab

**Parse cron expression:**

```go
ct, err := cron.Parse("0 0 0 */5 * *") // Schedule at every 5 days.
if err != nil {
    panic(err)
}
```

`NextAfter` return the next scheduled datetime after t.
`Next` returns the next scheduled datetime from now.

```go
// now = 2026-01-01 00:00:00 +0000 UTC
ct, _ := cron.Parse("0 0 0 */5 * *") // Schedule every 5 days.
ct.Next()                                                   // 2026-01-06 00:00:00 +0000 UTC
ct.NextAfter(time.Date(2026, 7, 31, 0, 0, 0, 0, time.UTC))  // 2026-08-01 00:00:00 +0000 UTC
ct.NextAfter(time.Date(2026, 8, 15, 0, 0, 0, 0, time.UTC))  // 2026-08-16 00:00:00 +0000 UTC
ct.NextAfter(time.Date(2026, 10, 01, 0, 0, 0, 0, time.UTC)) // 2026-10-06 00:00:00 +0000 UTC
```

### Crontab syntax

```txt
TZ=UTC * * * * * *
|      | | | | | |
|      | | | | | |- Day of week
|      | | | | |--- Month
|      | | | |----- Day of month
|      | | |------- Hour
|      | |--------- Minute
|      |----------- Second (Optional)
|------------------ Timezone (Optional)
```

| Field name   | Required |     Values      | Special characters |
| :----------- | :------: | :-------------: | :----------------: |
| Timezone     |    No    |  Timezone name  |                    |
| Second       |    No    |      0-59       |  `*` `/` `,` `-`   |
| Minute       |   Yes    |      0-59       |  `*` `/` `,` `-`   |
| Hours        |   Yes    |      0-23       |  `*` `/` `,` `-`   |
| Day of month |   Yes    |      1-31       |  `*` `/` `,` `-`   |
| Month        |   Yes    | 1-12 or JAN-DEC |  `*` `/` `,` `-`   |
| Day of week  |   Yes    | 0-6 or SUN-SAT  |  `*` `/` `,` `-`   |

> [!IMPORTANT]
> The "Day of month" and the "Day of week" are evaluated with AND condition.

Following aliases are defined for convenience.

| Alias name | Alias value | Usage example           |
| :--------- | :---------  | :---------------------: |
| CRON_TZ    | `TZ`        | `CRON_TZ=UTC 0 0 * * *` |
| @monthly   | `0 0 1 * *` | `TZ=UTC @monthly`       |
| @weekly    | `0 0 * * 0` | `TZ=UTC @weekly`        |
| @daily     | `0 0 * * *` | `TZ=UTC @daily`         |
| @hourly    | `0 * * * *` | `TZ=UTC @hourly`        |
| @sunday    | `0 0 * * 0` | `TZ=UTC @sunday`        |
| @monday    | `0 0 * * 1` | `TZ=UTC @monday`        |
| @tuesday   | `0 0 * * 2` | `TZ=UTC @tuesday`       |
| @wednesday | `0 0 * * 3` | `TZ=UTC @wednesday`     |
| @thursday  | `0 0 * * 4` | `TZ=UTC @thursday`      |
| @friday    | `0 0 * * 5` | `TZ=UTC @friday`        |
| @saturday  | `0 0 * * 6` | `TZ=UTC @saturday`      |

In addition, "@every <Duration>" expression can be used.
Format of the duration must follow the time.ParseDuration specifications.

Example of "@every" expression.

| Duration    | Resolved Cron          | Notes                 |
| :---------- | :--------------------- | :-------------------- |
| `-1s`       | `ERROR`                | Duration must be >0s  |
| `0s`        | `ERROR`                | Duration must be >0s  |
| `1s`        | `*/1 * * * * *`        |                       |
| `1m`        | `0 */1 * * * *`        |                       |
| `1h`        | `0 0 */1 * * *`        |                       |
| `61s`       | `*/1 */1 * * *`        |                       |
| `15m30s`    | `*/30 */15 * * * *`    |                       |
| `65m30s`    | `*/30 */5 */1 * * *`   |                       |
| `1h30m`     | `0 */30 */1 * * *`     |                       |
| `23h59m59s` | `*/59 */59 */23 * * *` |                       |
| `24h`       | `ERROR`                | Duration must be <24h |

### Run a job

Use `Cron` with job function.

```go
count := 0
c, err := cron.NewCron(&cron.Config{
    Crontab: "* * * * * *", // Every seconds.
    JobFunc: func(ctx context.Context) error {
        count++
        log.Println("Count:", count)
        return nil
    },
})
if err != nil {
    log.Panicln(err)
}

c.Start() // Start jobs to scheduled and run.
```

Cronjob is configureable.

```go
// Config is the configuration for the [Cron].
type Config struct {
    // Crontab is the cron expression.
    // See [Parse] for the syntax.
    Crontab string
    // MaxConcurrency is the maximum concurrency
    // of the currently running JobFunc.
    // Values should be 1, 2, ...
    // 1 means exactly 1 job at a time.
    // If less than 1, 1 is used.
    MaxConcurrency int
    // MaxRetry is the maximum number to retry
    // when the JobFunc returned non-nil error.
    // Values should be 0, 1, 2, ...
    // 0 means no retry. If less than 0, 0 is used.
    MaxRetry int
    // JobFunc is the job to run.
    JobFunc func(context.Context) error
    // WithContext provides a context which passed to the [Job.Run].
    // This is called once for a run and not called for retry.
    // [context.Background] is used when WithContext is nil.
    WithContext func() context.Context
    // EventHook is the function that hooks [Event]s.
    // The [Event] is notified through the first argument.
    // Additional information such as error is passed by a
    // depending on the event type.
    //  - For EventJobFailed: error is given by a[0].
    //  - For EventJobPanicked: recovered value given by a[0].
    EventHook func(e Event, a ...any)
}
```

## Docs & Examples

- GoDoc: <https://pkg.go.dev/github.com/aileron-projects/go-cron>
- Examples:
  - [example_test.go](./example_test.go)
  - Simple cron job: [examples/cronjob/](./examples/cronjob/)
  - Eventhook: [examples/eventhook/](./examples/eventhook/)

## References

- [robfig/cron](https://github.com/robfig/cron)
- [go-co-op/gocron](https://github.com/go-co-op/gocron)
- <https://en.wikipedia.org/wiki/Cron>
- <https://crontab.guru/>
- <https://crontab.cronhub.io/>
