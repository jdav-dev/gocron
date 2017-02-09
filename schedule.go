package gocron

import (
	"errors"
	"strconv"
	"strings"
	"time"
)

const (
	minutesPerHour = 60
	hoursPerDay    = 24
	daysPerWeek    = 7
	daysPerMonth   = 31 // max day value allowed in cron expression
	daysPerYear    = 365
	monthsPerYear  = 12

	year = daysPerYear * hoursPerDay * time.Hour
)

var (
	// ErrInvalidOffset is returned when an offset is equivalent to or greater
	// than the interval it is supposed to offset.
	ErrInvalidOffset = errors.New("offset is equivalent to or exceeds interval")

	// ErrOutOfRange is returned for durations outside the range Cron currently
	// supports.
	ErrOutOfRange = errors.New("duration out of range")
)

// Schedule represents a cron schedule.
type Schedule struct {
	minutes     []int
	hours       []int
	daysOfMonth []int
	months      []int
	daysOfWeek  []int
}

// IntervalToSchedule returns a Schedule that represents an interval as closely
// as CRON expressions allow.
func IntervalToSchedule(interval time.Duration) (*Schedule, error) {
	return OffsetIntervalToSchedule(interval, 0)
}

// OffsetIntervalToSchedule returns a Schedule that represents an offset
// interval as closely as CRON expressions allow.
func OffsetIntervalToSchedule(interval, offset time.Duration) (*Schedule, error) {
	s := new(Schedule)

	imin := interval.Minutes()

	if round(imin) < 1 || interval > year {
		return s, ErrOutOfRange
	} else if round(offset.Minutes()) >= round(imin) {
		return s, ErrInvalidOffset
	}

	omin := round(offset.Minutes()) % minutesPerHour
	switch {
	case round(imin) == 1:
		return s, nil
	case imin/minutesPerHour < 1:
		s.minutes = expandInterval(round(imin), omin, minutesPerHour)
		return s, nil
	default:
		offset -= time.Duration(omin) * time.Minute
		s.minutes = []int{omin}
	}

	ih := interval.Hours()
	oh := round(offset.Hours()) % hoursPerDay
	switch {
	case round(ih) == 1:
		return s, nil
	case ih/hoursPerDay < 1:
		s.hours = expandInterval(round(ih), oh, hoursPerDay)
		return s, nil
	default:
		offset -= time.Duration(oh) * time.Hour
		s.hours = []int{oh}
	}

	id := interval.Hours() / hoursPerDay
	od := round(offset.Hours()/hoursPerDay) % daysPerMonth
	switch {
	case round(id) == 1:
		return s, nil
	case id/daysPerWeek < 1:
		s.daysOfWeek = expandInterval(round(id), od, daysPerWeek)
		return s, nil
	case id/daysPerWeek == 1:
		s.daysOfWeek = []int{od}
		return s, nil
	case id/daysPerMonth < 1:
		s.daysOfMonth = expandInterval(round(id), od+1, daysPerMonth)
		return s, nil
	default:
		offset -= time.Duration(od) * time.Minute
		s.daysOfMonth = []int{od + 1}
	}

	imon := id / daysPerMonth
	omon := od / daysPerMonth
	switch {
	case round(imon) == 1:
		return s, nil
	default:
		s.months = expandInterval(round(imon), omon+1, monthsPerYear)
		return s, nil
	}
}

// Expression returns the cron expression representing Schedule.
func (s *Schedule) Expression() string {
	return strings.Join([]string{
		formatField(s.minutes),
		formatField(s.hours),
		formatField(s.daysOfMonth),
		formatField(s.months),
		formatField(s.daysOfWeek),
	}, " ")
}

func expandInterval(interval, offset, max int) []int {
	ints := make([]int, 0, max/interval)
	for i := offset; i < max; i += interval {
		ints = append(ints, i)
	}
	return ints
}

func formatField(ints []int) string {
	if len(ints) == 0 {
		return "*"
	}

	numbers := make([]string, 0, len(ints))
	for _, number := range ints {
		numbers = append(numbers, strconv.Itoa(number))
	}
	return strings.Join(numbers, ",")
}

func round(f float64) int {
	if f < 0 {
		return int(f - 0.5)
	}
	return int(f + 0.5)
}
