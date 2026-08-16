package cron

import (
	"testing"
	"time"

	"github.com/aileron-projects/go-tester"
)

func TestCrontab(t *testing.T) {
	t.Parallel()
	t.Run("nil time func", func(t *testing.T) {
		ct, err := Parse("* * * * *")
		tester.AssertEqual(t, nil, err)
		ct.WithTimeFunc(nil)
		ct.Next() // No panic.
	})
	t.Run("non nil time func", func(t *testing.T) {
		ct, err := Parse("* * * * *")
		tester.AssertEqual(t, nil, err)
		now := time.Now().Add(time.Hour)
		ct.WithTimeFunc(func() time.Time { return now })
		got := ct.Now()
		tester.AssertEqual(t, 0, now.Compare(got))
	})
}

func TestCrontab_NextAfter(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct {
		cron string
		t    time.Time
		want time.Time
	}{
		"case01":             {"TZ=Local */1 * * * * *", time.Date(2000, 1, 1, 0, 0, 0, 0, time.Local), time.Date(2000, 1, 1, 0, 0, 1, 0, time.Local)},
		"case02":             {"TZ=UTC   */1 * * * * *", time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC), time.Date(2000, 1, 1, 0, 0, 1, 0, time.UTC)},
		"case03":             {"TZ=UTC   */1 * * * * *", time.Date(2000, 1, 1, 0, 0, 59, 0, time.UTC), time.Date(2000, 1, 1, 0, 1, 0, 0, time.UTC)},     // Increment min.
		"case04":             {"TZ=UTC   */1 * * * * *", time.Date(2000, 1, 1, 0, 59, 59, 0, time.UTC), time.Date(2000, 1, 1, 1, 0, 0, 0, time.UTC)},    // Increment hour.
		"case05":             {"TZ=UTC   */1 * * * * *", time.Date(2000, 1, 1, 23, 59, 59, 0, time.UTC), time.Date(2000, 1, 2, 0, 0, 0, 0, time.UTC)},   // Increment day.
		"case06":             {"TZ=UTC   */1 * * * * *", time.Date(2000, 1, 31, 23, 59, 59, 0, time.UTC), time.Date(2000, 2, 1, 0, 0, 0, 0, time.UTC)},  // Increment month.
		"case07":             {"TZ=UTC   */1 * * * * *", time.Date(2000, 12, 31, 23, 59, 59, 0, time.UTC), time.Date(2001, 1, 1, 0, 0, 0, 0, time.UTC)}, // Increment year.
		"case08":             {"TZ=Local 0 */1 * * * *", time.Date(2000, 1, 1, 0, 0, 0, 0, time.Local), time.Date(2000, 1, 1, 0, 1, 0, 0, time.Local)},
		"case09":             {"TZ=UTC   0 */1 * * * *", time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC), time.Date(2000, 1, 1, 0, 1, 0, 0, time.UTC)},
		"case10":             {"TZ=UTC   0 */1 * * * *", time.Date(2000, 1, 1, 0, 59, 59, 0, time.UTC), time.Date(2000, 1, 1, 1, 0, 0, 0, time.UTC)},    // Increment hour.
		"case11":             {"TZ=UTC   0 */1 * * * *", time.Date(2000, 1, 1, 23, 59, 59, 0, time.UTC), time.Date(2000, 1, 2, 0, 0, 0, 0, time.UTC)},   // Increment day.
		"case12":             {"TZ=UTC   0 */1 * * * *", time.Date(2000, 1, 31, 23, 59, 59, 0, time.UTC), time.Date(2000, 2, 1, 0, 0, 0, 0, time.UTC)},  // Increment month.
		"case13":             {"TZ=UTC   0 */1 * * * *", time.Date(2000, 12, 31, 23, 59, 59, 0, time.UTC), time.Date(2001, 1, 1, 0, 0, 0, 0, time.UTC)}, // Increment year.
		"case14":             {"TZ=Local 0 0 */1 * * *", time.Date(2000, 1, 1, 0, 0, 0, 0, time.Local), time.Date(2000, 1, 1, 1, 0, 0, 0, time.Local)},
		"case15":             {"TZ=UTC   0 0 */1 * * *", time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC), time.Date(2000, 1, 1, 1, 0, 0, 0, time.UTC)},
		"case16":             {"TZ=UTC   0 0 */1 * * *", time.Date(2000, 1, 1, 23, 59, 59, 0, time.UTC), time.Date(2000, 1, 2, 0, 0, 0, 0, time.UTC)},   // Increment day.
		"case17":             {"TZ=UTC   0 0 */1 * * *", time.Date(2000, 1, 31, 23, 59, 59, 0, time.UTC), time.Date(2000, 2, 1, 0, 0, 0, 0, time.UTC)},  // Increment month.
		"case18":             {"TZ=UTC   0 0 */1 * * *", time.Date(2000, 12, 31, 23, 59, 59, 0, time.UTC), time.Date(2001, 1, 1, 0, 0, 0, 0, time.UTC)}, // Increment year.
		"case19":             {"TZ=Local 0 0 0 */1 * *", time.Date(2000, 1, 1, 0, 0, 0, 0, time.Local), time.Date(2000, 1, 2, 0, 0, 0, 0, time.Local)},
		"case20":             {"TZ=UTC   0 0 0 */1 * *", time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC), time.Date(2000, 1, 2, 0, 0, 0, 0, time.UTC)},
		"case21":             {"TZ=UTC   0 0 0 */1 * *", time.Date(2000, 1, 31, 23, 59, 59, 0, time.UTC), time.Date(2000, 2, 1, 0, 0, 0, 0, time.UTC)},  // Increment month.
		"case22":             {"TZ=UTC   0 0 0 */1 * *", time.Date(2000, 12, 31, 23, 59, 59, 0, time.UTC), time.Date(2001, 1, 1, 0, 0, 0, 0, time.UTC)}, // Increment year.
		"case23":             {"TZ=Local 0 0 0 1 */1 *", time.Date(2000, 1, 1, 0, 0, 0, 0, time.Local), time.Date(2000, 2, 1, 0, 0, 0, 0, time.Local)},
		"case24":             {"TZ=UTC   0 0 0 1 */1 *", time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC), time.Date(2000, 2, 1, 0, 0, 0, 0, time.UTC)},
		"case25":             {"TZ=UTC   0 0 0 1 */1 *", time.Date(2000, 12, 31, 23, 59, 59, 0, time.UTC), time.Date(2001, 1, 1, 0, 0, 0, 0, time.UTC)}, // Increment year.
		"case26":             {"TZ=Local */5 * * * * *", time.Date(2000, 1, 1, 0, 0, 0, 0, time.Local), time.Date(2000, 1, 1, 0, 0, 5, 0, time.Local)},
		"case27":             {"TZ=UTC   */5 * * * * *", time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC), time.Date(2000, 1, 1, 0, 0, 5, 0, time.UTC)},
		"case28":             {"TZ=UTC   */5 * * * * *", time.Date(2000, 1, 1, 0, 0, 59, 0, time.UTC), time.Date(2000, 1, 1, 0, 1, 0, 0, time.UTC)},     // Increment min.
		"case29":             {"TZ=UTC   */5 * * * * *", time.Date(2000, 1, 1, 0, 59, 59, 0, time.UTC), time.Date(2000, 1, 1, 1, 0, 0, 0, time.UTC)},    // Increment hour.
		"case30":             {"TZ=UTC   */5 * * * * *", time.Date(2000, 1, 1, 23, 59, 59, 0, time.UTC), time.Date(2000, 1, 2, 0, 0, 0, 0, time.UTC)},   // Increment day.
		"case31":             {"TZ=UTC   */5 * * * * *", time.Date(2000, 1, 31, 23, 59, 59, 0, time.UTC), time.Date(2000, 2, 1, 0, 0, 0, 0, time.UTC)},  // Increment month.
		"case32":             {"TZ=UTC   */5 * * * * *", time.Date(2000, 12, 31, 23, 59, 59, 0, time.UTC), time.Date(2001, 1, 1, 0, 0, 0, 0, time.UTC)}, // Increment year.
		"case33":             {"TZ=Local 0 */5 * * * *", time.Date(2000, 1, 1, 0, 0, 0, 0, time.Local), time.Date(2000, 1, 1, 0, 5, 0, 0, time.Local)},
		"case34":             {"TZ=UTC   0 */5 * * * *", time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC), time.Date(2000, 1, 1, 0, 5, 0, 0, time.UTC)},
		"case35":             {"TZ=UTC   0 */5 * * * *", time.Date(2000, 1, 1, 0, 59, 59, 0, time.UTC), time.Date(2000, 1, 1, 1, 0, 0, 0, time.UTC)},    // Increment hour.
		"case36":             {"TZ=UTC   0 */5 * * * *", time.Date(2000, 1, 1, 23, 59, 59, 0, time.UTC), time.Date(2000, 1, 2, 0, 0, 0, 0, time.UTC)},   // Increment day.
		"case37":             {"TZ=UTC   0 */5 * * * *", time.Date(2000, 1, 31, 23, 59, 59, 0, time.UTC), time.Date(2000, 2, 1, 0, 0, 0, 0, time.UTC)},  // Increment month.
		"case38":             {"TZ=UTC   0 */5 * * * *", time.Date(2000, 12, 31, 23, 59, 59, 0, time.UTC), time.Date(2001, 1, 1, 0, 0, 0, 0, time.UTC)}, // Increment year.
		"case39":             {"TZ=Local 0 0 */5 * * *", time.Date(2000, 1, 1, 0, 0, 0, 0, time.Local), time.Date(2000, 1, 1, 5, 0, 0, 0, time.Local)},
		"case40":             {"TZ=UTC   0 0 */5 * * *", time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC), time.Date(2000, 1, 1, 5, 0, 0, 0, time.UTC)},
		"case41":             {"TZ=UTC   0 0 */5 * * *", time.Date(2000, 1, 1, 23, 59, 59, 0, time.UTC), time.Date(2000, 1, 2, 0, 0, 0, 0, time.UTC)},   // Increment day.
		"case42":             {"TZ=UTC   0 0 */5 * * *", time.Date(2000, 1, 31, 23, 59, 59, 0, time.UTC), time.Date(2000, 2, 1, 0, 0, 0, 0, time.UTC)},  // Increment month.
		"case43":             {"TZ=UTC   0 0 */5 * * *", time.Date(2000, 12, 31, 23, 59, 59, 0, time.UTC), time.Date(2001, 1, 1, 0, 0, 0, 0, time.UTC)}, // Increment year.
		"case44":             {"TZ=Local 0 0 0 */5 * *", time.Date(2000, 1, 1, 0, 0, 0, 0, time.Local), time.Date(2000, 1, 6, 0, 0, 0, 0, time.Local)},
		"case45":             {"TZ=UTC   0 0 0 */5 * *", time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC), time.Date(2000, 1, 6, 0, 0, 0, 0, time.UTC)},
		"case46":             {"TZ=UTC   0 0 0 */5 * *", time.Date(2000, 1, 31, 23, 59, 59, 0, time.UTC), time.Date(2000, 2, 1, 0, 0, 0, 0, time.UTC)},  // Increment month.
		"case47":             {"TZ=UTC   0 0 0 */5 * *", time.Date(2000, 12, 31, 23, 59, 59, 0, time.UTC), time.Date(2001, 1, 1, 0, 0, 0, 0, time.UTC)}, // Increment year.
		"case48":             {"TZ=Local 0 0 0 1 */5 *", time.Date(2000, 1, 1, 0, 0, 0, 0, time.Local), time.Date(2000, 6, 1, 0, 0, 0, 0, time.Local)},
		"case49":             {"TZ=UTC   0 0 0 1 */5 *", time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC), time.Date(2000, 6, 1, 0, 0, 0, 0, time.UTC)},
		"case50":             {"TZ=UTC   0 0 0 1 */5 *", time.Date(2000, 12, 31, 23, 59, 59, 0, time.UTC), time.Date(2001, 1, 1, 0, 0, 0, 0, time.UTC)}, // Increment year.
		"case51":             {"TZ=Local 0 0 0 * * SUN", time.Date(2000, 1, 1, 0, 0, 0, 0, time.Local), time.Date(2000, 1, 2, 0, 0, 0, 0, time.Local)},  // 2000/01/01=SAT > 2000/01/01=SUN
		"case52":             {"TZ=UTC   0 0 0 * * SUN", time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC), time.Date(2000, 1, 2, 0, 0, 0, 0, time.UTC)},      // 2000/01/01=SAT > 2000/01/01=SUN
		"case53":             {"TZ=Local 0 0 0 * * FRI", time.Date(2000, 1, 1, 0, 0, 0, 0, time.Local), time.Date(2000, 1, 7, 0, 0, 0, 0, time.Local)},  // 2000/01/01=SAT > 2000/01/07=FRI
		"case54":             {"TZ=UTC   0 0 0 * * FRI", time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC), time.Date(2000, 1, 7, 0, 0, 0, 0, time.UTC)},      // 2000/01/01=SAT > 2000/01/07=FRI
		"Day31_01":           {"0 0 0 31 * *", time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC), time.Date(2000, 1, 31, 0, 0, 0, 0, time.UTC)},               // 1/1 > 1/31
		"Day31_02":           {"0 0 0 31 * *", time.Date(2000, 2, 1, 0, 0, 0, 0, time.UTC), time.Date(2000, 3, 31, 0, 0, 0, 0, time.UTC)},               // 2/1 > 3/31
		"Day31_03":           {"0 0 0 31 * *", time.Date(2000, 12, 31, 0, 0, 0, 0, time.UTC), time.Date(2001, 1, 31, 0, 0, 0, 0, time.UTC)},             // 12/31 > 1/31
		"LeapYear_01":        {"0 0 0 29 2 *", time.Date(2000, 2, 1, 0, 0, 0, 0, time.UTC), time.Date(2000, 2, 29, 0, 0, 0, 0, time.UTC)},               // 2000/02/01 > 2000/02/29
		"LeapYear_02":        {"0 0 0 29 2 *", time.Date(2000, 2, 29, 0, 0, 0, 0, time.UTC), time.Date(2004, 2, 29, 0, 0, 0, 0, time.UTC)},              // 2000/02/29 > 2004/02/29
		"LeapYear_03":        {"0 0 0 29 2 *", time.Date(2000, 3, 1, 0, 0, 0, 0, time.UTC), time.Date(2004, 2, 29, 0, 0, 0, 0, time.UTC)},               // 2000/03/01 > 2004/02/29
		"Day&WeekDay_01":     {"0 0 0 23 * SUN", time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC), time.Date(2000, 1, 23, 0, 0, 0, 0, time.UTC)},             // 2000/01/01=SUN > 2 2000/01/23=SUN
		"Day&WeekDay_02":     {"0 0 0 23 * SUN", time.Date(2000, 1, 23, 0, 0, 0, 0, time.UTC), time.Date(2000, 4, 23, 0, 0, 0, 0, time.UTC)},            // 2000/01/23=SUN > 2000/04/23=SUN
		"At 1sec":            {"1 * * * * *", time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC), time.Date(2000, 1, 1, 0, 0, 1, 0, time.UTC)},
		"At 1min":            {"1 1 * * * *", time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC), time.Date(2000, 1, 1, 0, 1, 1, 0, time.UTC)},
		"At 1o'clock":        {"1 1 1 * * *", time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC), time.Date(2000, 1, 1, 1, 1, 1, 0, time.UTC)},
		"range sec 01":       {"10-30,45 * * * * *", time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC), time.Date(2000, 1, 1, 0, 0, 10, 0, time.UTC)},      // 00:00:00 > 00:00:10
		"range sec 02":       {"10-30,45 * * * * *", time.Date(2000, 1, 1, 0, 0, 10, 0, time.UTC), time.Date(2000, 1, 1, 0, 0, 11, 0, time.UTC)},     // 00:00:10 > 00:00:11
		"range sec 03":       {"10-30,45 * * * * *", time.Date(2000, 1, 1, 0, 0, 30, 0, time.UTC), time.Date(2000, 1, 1, 0, 0, 45, 0, time.UTC)},     // 00:00:30 > 00:00:45
		"range min 01":       {"0 10-30,45 * * * *", time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC), time.Date(2000, 1, 1, 0, 10, 0, 0, time.UTC)},      // 00:00:00 > 00:10:00
		"range min 02":       {"0 10-30,45 * * * *", time.Date(2000, 1, 1, 0, 10, 0, 0, time.UTC), time.Date(2000, 1, 1, 0, 11, 00, 0, time.UTC)},    // 00:10:00 > 00:11:00
		"range min 03":       {"0 10-30,45 * * * *", time.Date(2000, 1, 1, 0, 30, 0, 0, time.UTC), time.Date(2000, 1, 1, 0, 45, 00, 0, time.UTC)},    // 00:30:00 > 00:45:00
		"range hour 01":      {"0 0 15-20,23 * * *", time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC), time.Date(2000, 1, 1, 15, 0, 0, 0, time.UTC)},      // 00:00:00 > 15:00:00
		"range hour 02":      {"0 0 15-20,23 * * *", time.Date(2000, 1, 1, 15, 0, 0, 0, time.UTC), time.Date(2000, 1, 1, 16, 00, 00, 0, time.UTC)},   // 15:00:00 > 16:00:00
		"range hour 03":      {"0 0 15-20,23 * * *", time.Date(2000, 1, 1, 20, 0, 0, 0, time.UTC), time.Date(2000, 1, 1, 23, 00, 00, 0, time.UTC)},   // 20:00:00 > 23:00:00
		"range step sec 01":  {"10-30/3,45 * * * * *", time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC), time.Date(2000, 1, 1, 0, 0, 10, 0, time.UTC)},    // 00:00:00 > 00:00:10
		"range step sec 02":  {"10-30/3,45 * * * * *", time.Date(2000, 1, 1, 0, 0, 10, 0, time.UTC), time.Date(2000, 1, 1, 0, 0, 13, 0, time.UTC)},   // 00:00:10 > 00:00:13
		"range step sec 03":  {"10-30/3,45 * * * * *", time.Date(2000, 1, 1, 0, 0, 30, 0, time.UTC), time.Date(2000, 1, 1, 0, 0, 45, 0, time.UTC)},   // 00:00:30 > 00:00:45
		"range step min 01":  {"0 10-30/3,45 * * * *", time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC), time.Date(2000, 1, 1, 0, 10, 0, 0, time.UTC)},    // 00:00:00 > 00:10:00
		"range step min 02":  {"0 10-30/3,45 * * * *", time.Date(2000, 1, 1, 0, 10, 0, 0, time.UTC), time.Date(2000, 1, 1, 0, 13, 00, 0, time.UTC)},  // 00:10:00 > 00:13:00
		"range step min 03":  {"0 10-30/3,45 * * * *", time.Date(2000, 1, 1, 0, 30, 0, 0, time.UTC), time.Date(2000, 1, 1, 0, 45, 00, 0, time.UTC)},  // 00:30:00 > 00:45:00
		"range step hour 01": {"0 0 15-20/3,23 * * *", time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC), time.Date(2000, 1, 1, 15, 0, 0, 0, time.UTC)},    // 00:00:00 > 15:00:00
		"range step hour 02": {"0 0 15-20/3,23 * * *", time.Date(2000, 1, 1, 15, 0, 0, 0, time.UTC), time.Date(2000, 1, 1, 18, 00, 00, 0, time.UTC)}, // 15:00:00 > 18:00:00
		"range step hour 03": {"0 0 15-20/3,23 * * *", time.Date(2000, 1, 1, 20, 0, 0, 0, time.UTC), time.Date(2000, 1, 1, 23, 00, 00, 0, time.UTC)}, // 20:00:00 > 23:00:00
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			ct, err := Parse(tc.cron)
			tester.AssertEqual(t, nil, err)
			got := ct.NextAfter(tc.t)
			tester.AssertEqual(t, 0, tc.want.Compare(got))
			t.Log(tc.want, got)
		})
	}
}

func TestNextTime(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct {
		targets uint64
		now     int
		max     int
		want    int
	}{
		"case00": {0, 0, 9, 0},
		"case01": {0b_00011, 0, 9, 1},
		"case02": {0b_00101, 0, 9, 2},
		"case03": {0b_10000_00001, 0, 9, 9},
		"case04": {0b_00001_00000_00001, 0, 9, 0},
	}
	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			got := nextTime(tc.targets, tc.now, tc.max)
			tester.AssertEqual(t, tc.want, got)
		})
	}
}

func TestParse(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct {
		exp string
		err error
		ct  *Crontab
	}{
		"case01": {
			"* * * * *", nil,
			&Crontab{
				second: 0b_00001,                                                                   // 0
				minute: 0b_11111_11111_11111_11111_11111_11111_11111_11111_11111_11111_11111_11111, // 0-59
				hour:   0b_01111_11111_11111_11111_11111,                                           // 0-23
				day:    0b_00011_11111_11111_11111_11111_11111_11110,                               // 1-31
				month:  0b_00111_11111_11110,                                                       // 1-12
				week:   0b_00011_11111,                                                             // 0-6
			},
		},
		"case02": {
			"* * * * * *", nil,
			&Crontab{
				second: 0b_11111_11111_11111_11111_11111_11111_11111_11111_11111_11111_11111_11111,
				minute: 0b_11111_11111_11111_11111_11111_11111_11111_11111_11111_11111_11111_11111,
				hour:   0b_01111_11111_11111_11111_11111,
				day:    0b_00011_11111_11111_11111_11111_11111_11110,
				month:  0b_00111_11111_11110,
				week:   0b_00011_11111,
			},
		},
		"case03": {
			"TZ=UTC * * * * *", nil,
			&Crontab{
				second: 0b_00001,                                                                   // 0
				minute: 0b_11111_11111_11111_11111_11111_11111_11111_11111_11111_11111_11111_11111, // 0-59
				hour:   0b_01111_11111_11111_11111_11111,                                           // 0-23
				day:    0b_00011_11111_11111_11111_11111_11111_11110,                               // 1-31
				month:  0b_00111_11111_11110,                                                       // 1-12
				week:   0b_00011_11111,                                                             // 0-6
			},
		},
		"case04": {
			"CRON_TZ=UTC * * * * *", nil,
			&Crontab{
				second: 0b_00001,                                                                   // 0
				minute: 0b_11111_11111_11111_11111_11111_11111_11111_11111_11111_11111_11111_11111, // 0-59
				hour:   0b_01111_11111_11111_11111_11111,                                           // 0-23
				day:    0b_00011_11111_11111_11111_11111_11111_11110,                               // 1-31
				month:  0b_00111_11111_11110,                                                       // 1-12
				week:   0b_00011_11111,                                                             // 0-6
			},
		},
		"case05": {
			"TZ=UTC * * * * * *", nil,
			&Crontab{
				second: 0b_11111_11111_11111_11111_11111_11111_11111_11111_11111_11111_11111_11111,
				minute: 0b_11111_11111_11111_11111_11111_11111_11111_11111_11111_11111_11111_11111,
				hour:   0b_01111_11111_11111_11111_11111,
				day:    0b_00011_11111_11111_11111_11111_11111_11110,
				month:  0b_00111_11111_11110,
				week:   0b_00011_11111,
			},
		},
		"case06": {
			"CRON_TZ=UTC * * * * * *", nil,
			&Crontab{
				second: 0b_11111_11111_11111_11111_11111_11111_11111_11111_11111_11111_11111_11111,
				minute: 0b_11111_11111_11111_11111_11111_11111_11111_11111_11111_11111_11111_11111,
				hour:   0b_01111_11111_11111_11111_11111,
				day:    0b_00011_11111_11111_11111_11111_11111_11110,
				month:  0b_00111_11111_11110,
				week:   0b_00011_11111,
			},
		},
		"case07": {
			"0 0 1 1 *", nil,
			&Crontab{
				second: 0b_00001,
				minute: 0b_00001,
				hour:   0b_00001,
				day:    0b_00010,
				month:  0b_00010,
				week:   0b_00011_11111,
			},
		},
		"case08": {
			"0 0 0 1 1 *", nil,
			&Crontab{
				second: 0b_00001,
				minute: 0b_00001,
				hour:   0b_00001,
				day:    0b_00010,
				month:  0b_00010,
				week:   0b_00011_11111,
			},
		},
		"invalid second": {
			"x * * * * *",
			&ParseError{What: "second"},
			nil,
		},
		"invalid minute": {
			"* x * * * *",
			&ParseError{What: "minute"},
			nil,
		},
		"invalid hour": {
			"* * x * * *",
			&ParseError{What: "hour"},
			nil,
		},
		"invalid day of month": {
			"* * * x * *",
			&ParseError{What: "day of month"},
			nil,
		},
		"invalid month": {
			"* * * * x *",
			&ParseError{What: "month"},
			nil,
		},
		"invalid day of week": {
			"* * * * * x",
			&ParseError{What: "day of week"},
			nil,
		},
		"invalid number of fields": {
			"* * * * * * *",
			&ParseError{What: "number of fields"},
			nil,
		},
		"invalid location": {
			"TZ=NotExists * * * * *",
			&ParseError{What: "timezone"},
			nil,
		},
		"unschedulable": {
			"* * 30 2 *", // Feb 30th
			&ParseError{What: "schedule (unschedulable)"},
			nil,
		},
		"empty cron": {
			"",
			&ParseError{What: "cron expression"},
			nil,
		},
	}
	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			ct, err := Parse(tc.exp)
			tester.AssertEqualErr(t, tc.err, err)
			if tc.err != nil {
				tester.AssertEqual(t, nil, ct)
				return
			}
			tester.AssertEqual(t, tc.ct.second, ct.second)
			tester.AssertEqual(t, tc.ct.minute, ct.minute)
			tester.AssertEqual(t, tc.ct.hour, ct.hour)
			tester.AssertEqual(t, tc.ct.day, ct.day)
			tester.AssertEqual(t, tc.ct.month, ct.month)
			tester.AssertEqual(t, tc.ct.week, ct.week)
		})
	}
}

func TestReplaceAlias(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct {
		exp  string
		want string
	}{
		"case01": {"@monthly", "0 0 1 * *"},
		"case02": {"@MONTHLY", "0 0 1 * *"},
		"case03": {"CRON_TZ=UTC @monthly", "TZ=UTC 0 0 1 * *"},
		"case04": {"@weekly", "0 0 * * 0"},
		"case05": {"@WEEKLY", "0 0 * * 0"},
		"case06": {"CRON_TZ=UTC @weekly", "TZ=UTC 0 0 * * 0"},
		"case07": {"@daily", "0 0 * * *"},
		"case08": {"@DAILY", "0 0 * * *"},
		"case09": {"CRON_TZ=UTC @daily", "TZ=UTC 0 0 * * *"},
		"case10": {"@hourly", "0 * * * *"},
		"case12": {"@HOURLY", "0 * * * *"},
		"case11": {"CRON_TZ=UTC @hourly", "TZ=UTC 0 * * * *"},
		"case13": {"@sunday", "0 0 * * 0"},
		"case14": {"@monday", "0 0 * * 1"},
		"case15": {"@tuesday", "0 0 * * 2"},
		"case16": {"@thursday", "0 0 * * 4"},
		"case17": {"@friday", "0 0 * * 5"},
		"case18": {"@saturday", "0 0 * * 6"},
		"case19": {"@@", "@@"},
		"case20": {"@foo", "@foo"},
		"case21": {"FOO", "FOO"},
		"case22": {"X", "X"},
		"case23": {"@every", "@every"},
		"case24": {"@every -1s", "@every -1s"},
		"case25": {"@every 0s", "@every 0s"},
		"case26": {"@every 24h", "@every 24h"},
		"case27": {"@every 23h", "0 0 */23 * * *"},
		"case28": {"@every 60m", "0 0 */1 * * *"},
		"case29": {"@every 59m", "0 */59 * * * *"},
		"case30": {"@every 60s", "0 */1 * * * *"},
		"case31": {"@every 59s", "*/59 * * * * *"},
		"case32": {"@every 30s", "*/30 * * * * *"},
		"case33": {"@every 1s", "*/1 * * * * *"},
		"case34": {"@every 23h59m59s", "*/59 */59 */23 * * *"},
		"case35": {"@every 61m1s", "*/1 */1 */1 * * *"},
		"case36": {"@every 61m61s", "*/1 */2 */1 * * *"},
		"case37": {"@every 61s", "*/1 */1 * * * *"},
		"case38": {"TZ=UTC @every 61s", "TZ=UTC */1 */1 * * * *"},
	}
	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			got := replaceAlias(tc.exp)
			tester.AssertEqual(t, tc.want, got)
		})
	}
}

func TestNormalizeMonth(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct {
		exp  string
		want string
	}{
		"case01": {"JAN,FEB,MAR,APR,MAY,JUN,JUL,AUG,SEP,OCT,NOV,DEC", "1,2,3,4,5,6,7,8,9,10,11,12"},
		"case02": {"jan,feb,mar,apr,may,jun,jul,aug,sep,oct,nov,dec", "1,2,3,4,5,6,7,8,9,10,11,12"},
		"case03": {"Jan,Feb,Mar,Apr,May,Jun,Jul,Aug,Sep,Oct,Nov,Dec", "1,2,3,4,5,6,7,8,9,10,11,12"},
		"case04": {"JAN-DEC", "1-12"},
		"case05": {"JANFEB", "JANFEB"},
		"case06": {"FOO", "FOO"},
		"case07": {"foo", "foo"},
		"case08": {"x", "x"},
	}
	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			got := normalizeMonth(tc.exp)
			tester.AssertEqual(t, tc.want, got)
		})
	}
}

func TestNormalizeWeek(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct {
		exp  string
		want string
	}{
		"case01": {"SUN,MON,TUE,WED,THU,FRI,SAT", "0,1,2,3,4,5,6"},
		"case02": {"sun,mon,tue,wed,thu,fri,sat", "0,1,2,3,4,5,6"},
		"case03": {"Sun,Mon,Tue,Wed,Thu,Fri,Sat", "0,1,2,3,4,5,6"},
		"case04": {"SUN-SAT", "0-6"},
		"case05": {"SUNMON", "SUNMON"},
		"case06": {"FOO", "FOO"},
		"case07": {"foo", "foo"},
		"case08": {"x", "x"},
	}
	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			got := normalizeWeek(tc.exp)
			tester.AssertEqual(t, tc.want, got)
		})
	}
}

func TestParseValue(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct {
		exp      string // args
		min, max int    // args
		cron     uint64 // want
		err      string // want
	}{
		"case01": {"*", -1, -1, 0, "require 0 <= min <= max"},
		"case02": {"*", -1, 1, 0, "require 0 <= min <= max"},
		"case03": {"*", 1, -1, 0, "require 0 <= min <= max"},
		"case04": {"*", 2, 1, 0, "require 0 <= min <= max"},
		"case05": {"*", 0, 9, 0b_11111_11111, ""},
		"case06": {"*", 3, 9, 0b_11111_11000, ""},
		"case07": {"5", 0, 9, 0b_00001_00000, ""},
		"case08": {"5", 3, 9, 0b_00001_00000, ""},
		"case09": {"5", 6, 9, 0, "value should be >=6: `5`"},
		"case10": {"5", 3, 4, 0, "value should be <=4: `5`"},
		"case11": {"5-7", 0, 9, 0b_00111_00000, ""},
		"case12": {"5-7", 5, 9, 0b_00111_00000, ""},
		"case13": {"5-7", 6, 9, 0, "range should be 6<=start<=end<=9: 5-7"},
		"case14": {"5-7", 0, 6, 0, "range should be 0<=start<=end<=6: 5-7"},
		"case15": {"*/2", 0, 9, 0b_01010_10101, ""},
		"case16": {"*/3", 0, 9, 0b_10010_01001, ""},
		"case17": {"1/3", 0, 9, 0b_00100_10010, ""},
		"case18": {"10/3", 0, 9, 0, "value should be <=9: `10`"}, // Our of range,
		"case19": {"1-6/2", 0, 9, 0b_00001_01010, ""},
		"case20": {"6-1/2", 0, 9, 0, "range should be 0<=start<=end<=9: 6-1"},
		"case21": {"x-6/2", 0, 9, 0, "not an integer: `x`"},
		"case22": {"1-x/2", 0, 9, 0, "not an integer: `x`"},
		"case23": {"1-5/x", 0, 9, 0, "not an integer: `x`"},
		"case24": {"1-5/0", 0, 9, 0, "step should be 1<=step<=10: `0`"},
		"case25": {"5/99", 0, 9, 0, "step should be 1<=step<=10: `99`"},
		"case26": {"5/1", 0, 9, 0b_11111_00000, ""},
		"case27": {"5/9", 0, 9, 0b_00001_00000, ""},
		"case28": {"5/10", 0, 9, 0b_00001_00000, ""},
		"case29": {"1,2,5", 0, 9, 0b_00001_00110, ""},
		"case2:": {"1,2,5,*", 0, 9, 0b_11111_11111, ""},
		"case2;": {"1,2,5-7", 0, 9, 0b_00111_00110, ""},
		"case2<": {"*/2,5-7", 0, 9, 0b_01111_10101, ""},
		"case2=": {"*/2,5-7,9", 0, 9, 0b_11111_10101, ""},
		"case2>": {"1/2/3", 0, 9, 0, "invalid expression: 1/2/3"},
	}
	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			cron, err := parseValue(tc.exp, tc.min, tc.max)
			tester.AssertEqual(t, tc.cron, cron)
			if err != nil {
				tester.AssertEqual(t, tc.err, err.Error())
			} else {
				tester.AssertEqual(t, "", tc.err)
			}
		})
	}
}

func TestParseRange(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct {
		exp          string // args
		min, max     int    // args
		under, upper int    // want
		err          string // want
	}{
		"case01": {"*", -1, -1, 0, 0, "require 0 <= min <= max"},
		"case02": {"*", -1, 1, 0, 0, "require 0 <= min <= max"},
		"case03": {"*", 1, -1, 0, 0, "require 0 <= min <= max"},
		"case04": {"*", 2, 1, 0, 0, "require 0 <= min <= max"},
		"case05": {"*", 1, 2, 1, 2, ""},
		"case06": {"*", 1, 10, 1, 10, ""},
		"case07": {"**", 1, 10, 0, 0, "not an integer: `**`"}, // Invalid wildcard.
		"case08": {"5", 1, 10, 5, 5, ""},
		"case09": {"5", 5, 10, 5, 5, ""},
		"case10": {"5", 1, 5, 5, 5, ""},
		"case11": {"5", 8, 10, 0, 0, "value should be >=8: `5`"}, // Out of range.
		"case12": {"5x", 1, 10, 0, 0, "not an integer: `5x`"},    // Invalid number.
		"case13": {"x", 1, 10, 0, 0, "not an integer: `x`"},      // Not a number.
		"case14": {"2-8", 1, 10, 2, 8, ""},
		"case15": {"2-8", 1, 8, 2, 8, ""},  // Boundary value check.
		"case16": {"2-8", 2, 10, 2, 8, ""}, // Boundary value check.
		"case17": {"5-5", 1, 10, 5, 5, ""},
		"case18": {"2-8", 3, 10, 0, 0, "range should be 3<=start<=end<=10: 2-8"}, // Min out of range.
		"case19": {"2-8", 1, 5, 0, 0, "range should be 1<=start<=end<=5: 2-8"},   // Max out of range.
		"case20": {"8-2", 1, 10, 0, 0, "range should be 1<=start<=end<=10: 8-2"}, // Min>Max.
		"case21": {"2x-8", 1, 10, 0, 0, "not an integer: `2x`"},                  // Invalid range expression.
		"case22": {"2-5-8", 1, 10, 0, 0, "invalid expression: 2-5-8"},            // Invalid range expression.
		"case23": {"2-", 1, 10, 0, 0, "not an integer: ``"},                      // Invalid range expression.
		"case24": {"-8", 1, 10, 0, 0, "not an integer: ``"},                      // Invalid range expression.
	}
	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			under, upper, err := parseRange(tc.exp, tc.min, tc.max)
			tester.AssertEqual(t, tc.under, under)
			tester.AssertEqual(t, tc.upper, upper)
			if err != nil {
				tester.AssertEqual(t, tc.err, err.Error())
			} else {
				tester.AssertEqual(t, "", tc.err)
			}
		})
	}
}
