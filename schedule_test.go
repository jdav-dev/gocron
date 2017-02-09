package gocron

import (
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

type scheduleTest struct {
	interval time.Duration
	offset   time.Duration
	expected string
}

var (
	scheduleTests = []scheduleTest{
		scheduleTest{
			interval: time.Minute / 2,
			expected: "* * * * *",
		},
		scheduleTest{
			interval: time.Minute,
			expected: "* * * * *",
		},
		scheduleTest{
			interval: 15 * time.Minute,
			expected: "0,15,30,45 * * * *",
		},
		scheduleTest{
			interval: 15 * time.Minute,
			offset:   time.Minute,
			expected: "1,16,31,46 * * * *",
		},
		scheduleTest{
			interval: 15 * time.Minute,
			offset:   5 * time.Minute,
			expected: "5,20,35,50 * * * *",
		},
		scheduleTest{
			interval: time.Hour,
			expected: "0 * * * *",
		},
		scheduleTest{
			interval: time.Hour + 20*time.Minute,
			expected: "0 * * * *",
		},
		scheduleTest{
			interval: time.Hour,
			offset:   15 * time.Minute,
			expected: "15 * * * *",
		},
		scheduleTest{
			interval: 13 * time.Hour,
			expected: "0 0,13 * * *",
		},
		scheduleTest{
			interval: 24 * time.Hour,
			expected: "0 0 * * *",
		},
		scheduleTest{
			interval: 24 * time.Hour,
			offset:   time.Minute,
			expected: "1 0 * * *",
		},
		scheduleTest{
			interval: 24 * time.Hour,
			offset:   time.Hour,
			expected: "0 1 * * *",
		},
		scheduleTest{
			interval: 24 * time.Hour,
			offset:   time.Hour + 15*time.Minute,
			expected: "0 1 * * *",
		},
		scheduleTest{
			interval: 7 * hoursPerDay * time.Hour, // 1 week
			expected: "0 0 * * 0",
		},
		scheduleTest{
			interval: daysPerMonth * hoursPerDay * time.Hour, // ~1 month
			expected: "0 0 1 * *",
		},
		scheduleTest{
			interval: 800 * time.Hour,
			expected: "0 0 1 * *",
		},
		scheduleTest{
			interval: 4 * daysPerMonth * hoursPerDay * time.Hour, // ~4 months
			expected: "0 0 1 1,5,9 *",
		},
		scheduleTest{
			interval: 6 * daysPerMonth * hoursPerDay * time.Hour, // ~6 months
			expected: "0 0 1 1,7 *",
		},
		scheduleTest{
			interval: 6 * daysPerMonth * hoursPerDay * time.Hour, // ~6 months
			offset:   15 * hoursPerDay * time.Hour,               // 15 days
			expected: "0 0 16 1,7 *",
		},
		scheduleTest{
			interval: 365 * hoursPerDay * time.Hour, // 1 year
			expected: "0 0 1 1 *",
		},
		scheduleTest{
			interval: 365 * hoursPerDay * time.Hour,
			offset:   20 * time.Minute,
			expected: "20 0 1 1 *",
		},
	}
)

func TestOffsetIntervalToSchedule(t *testing.T) {
	for _, test := range scheduleTests {
		runTest(t, test)
	}
}

func runTest(t *testing.T, test scheduleTest) {
	s, err := OffsetIntervalToSchedule(test.interval, test.offset)
	assert.NoError(t, err)
	message := fmt.Sprintf("interval: %s, offset: %s", test.interval, test.offset)
	assert.Equal(t, test.expected, s.Expression(), message)
}

func TestErrOutOfRange(t *testing.T) {
	assert := assert.New(t)

	halfMin := time.Minute / 2

	_, err := IntervalToSchedule(halfMin)
	assert.NoError(err, halfMin.String())

	_, err = IntervalToSchedule(halfMin - time.Nanosecond)
	assert.Equal(ErrOutOfRange, err)

	_, err = IntervalToSchedule(year)
	assert.NoError(err)

	_, err = IntervalToSchedule(year + time.Nanosecond)
	assert.Equal(ErrOutOfRange, err)
}

func TestErrInvalidOffset(t *testing.T) {
	assert := assert.New(t)

	var interval, offset time.Duration

	interval, offset = time.Minute, 0
	_, err := OffsetIntervalToSchedule(interval, offset)
	message := fmt.Sprintf("interval: %s, offset: %s", interval, offset)
	assert.NoError(err, message)

	interval, offset = time.Hour, time.Hour/2
	_, err = OffsetIntervalToSchedule(interval, offset)
	message = fmt.Sprintf("interval: %s, offset: %s", interval, offset)
	assert.NoError(err, message)

	interval, offset = time.Minute, time.Minute
	_, err = OffsetIntervalToSchedule(interval, offset)
	message = fmt.Sprintf("interval: %s, offset: %s", interval, offset)
	assert.Equal(ErrInvalidOffset, err, message)

	interval, offset = time.Hour, 2*time.Hour
	_, err = OffsetIntervalToSchedule(interval, offset)
	message = fmt.Sprintf("interval: %s, offset: %s", interval, offset)
	assert.Equal(ErrInvalidOffset, err, message)
}
