package klog

import (
	"cloud.google.com/go/civil"
	"errors"
	"fmt"
	"regexp"
	"strconv"
	gotime "time"
)

// Time represents a wall clock time. It can be shifted to the adjacent dates.
type Time interface {
	Hour() int
	Minute() int

	// MidnightOffset returns the duration since (positive) or until (negative) midnight.
	MidnightOffset() Duration

	// IsYesterday checks whether the time is shifted to the previous day.
	IsYesterday() bool

	// IsTomorrow checks whether the time is shifted to the next day.
	IsTomorrow() bool

	// IsToday checks whether the time is not shifted.
	IsToday() bool
	IsEqualTo(Time) bool
	IsAfterOrEqual(Time) bool

	// Plus returns a time, where the specified duration was added. It doesn’t modify
	// the original object. If the resulting time would be shifted by more than one
	// day, it returns an error.
	Plus(Duration) (Time, error)

	// ToString serialises the time, e.g. `8:00` or `23:00>`
	ToString() string

	// ToStringWithFormat serialises the date according to the given format.
	ToStringWithFormat(TimeFormat) string

	// Format returns the current formatting.
	Format() TimeFormat
}

// TimeFormat contains the formatting options for the Time.
type TimeFormat struct {
	Use24HourClock bool
}

// DefaultTimeFormat returns the canonical time format, as recommended by the spec.
func DefaultTimeFormat() TimeFormat {
	return TimeFormat{
		Use24HourClock: true,
	}
}

type time struct {
	hour     int
	minute   int
	dayShift int
	format   TimeFormat
}

func newTime(hour int, minute int, dayShift int, format TimeFormat) (Time, error) {
	if hour == 24 && minute == 00 && dayShift <= 0 {
		// Accept a time of 24:00 (today), and interpret it as 0:00 (tomorrow).
		// Accept a time of 24:00 (yesterday), and interpret it as 0:00 (today).
		// This case is not supported for 24:00 (tomorrow), since that couldn’t be represented.
		hour = 0
		dayShift += 1
	}
	ct := civil.Time{Hour: hour, Minute: minute}
	if !ct.IsValid() {
		return nil, errors.New("INVALID_TIME")
	}
	return &time{
		hour:     ct.Hour,
		minute:   ct.Minute,
		dayShift: dayShift,
		format:   format,
	}, nil
}

func NewTime(hour int, minute int) (Time, error) {
	return newTime(hour, minute, 0, DefaultTimeFormat())
}

func NewTimeYesterday(hour int, minute int) (Time, error) {
	return newTime(hour, minute, -1, DefaultTimeFormat())
}

func NewTimeTomorrow(hour int, minute int) (Time, error) {
	return newTime(hour, minute, +1, DefaultTimeFormat())
}

var timePattern = regexp.MustCompile(`^(<)?(\d{1,2}):(\d{2})(am|pm)?(>)?$`)

func NewTimeFromString(hhmm string) (Time, error) {
	match := timePattern.FindStringSubmatch(hhmm)
	if len(match) != 6 || (match[1] == "<" && match[5] == ">") {
		return nil, errors.New("MALFORMED_TIME")
	}
	hour, _ := strconv.Atoi(match[2])
	minute, _ := strconv.Atoi(match[3])
	format := DefaultTimeFormat()
	if match[4] == "am" || match[4] == "pm" {
		if hour < 1 || hour > 12 {
			return nil, errors.New("INVALID_TIME")
		}
		format.Use24HourClock = false
		if match[4] == "am" && hour == 12 {
			hour = 0
		} else if match[4] == "pm" && hour < 12 {
			hour += 12
		}
	}
	dayShift := 0
	if match[1] == "<" {
		dayShift = -1
	} else if match[5] == ">" {
		dayShift = +1
	}
	return newTime(hour, minute, dayShift, format)
}

func NewTimeFromGo(t gotime.Time) Time {
	time, err := NewTime(t.Hour(), t.Minute())
	if err != nil {
		// This can/should never occur
		panic("ILLEGAL_TIME")
	}
	return time
}

func (t *time) Hour() int {
	return t.hour
}

func (t *time) Minute() int {
	return t.minute
}

func (t *time) MidnightOffset() Duration {
	if t.IsYesterday() {
		return NewDuration(-23+t.Hour(), -60+t.Minute())
	} else if t.IsTomorrow() {
		return NewDuration(24+t.Hour(), t.Minute())
	}
	return NewDuration(t.Hour(), t.Minute())
}

func (t *time) IsToday() bool {
	return t.dayShift == 0
}

func (t *time) IsYesterday() bool {
	return t.dayShift < 0
}

func (t *time) IsTomorrow() bool {
	return t.dayShift > 0
}

func (t *time) IsEqualTo(otherTime Time) bool {
	return t.MidnightOffset().InMinutes() == otherTime.MidnightOffset().InMinutes()
}

func (t *time) IsAfterOrEqual(otherTime Time) bool {
	first := t.MidnightOffset()
	second := otherTime.MidnightOffset()
	return first.InMinutes() >= second.InMinutes()
}

func (t *time) Plus(d Duration) (Time, error) {
	ONE_DAY := 24 * 60
	mins := t.MidnightOffset().Plus(d).InMinutes()
	if mins >= 2*ONE_DAY || mins < ONE_DAY*-1 {
		return nil, errors.New("IMPOSSIBLE_OPERATION")
	}
	dayShift := 0
	if mins < 0 {
		dayShift = -1
		mins = ONE_DAY + mins
	} else if mins > ONE_DAY {
		dayShift = 1
		mins = mins - ONE_DAY
	}
	return newTime(mins/60, mins%60, dayShift, t.format)
}

func (t *time) ToString() string {
	yesterdayPrefix := ""
	if t.IsYesterday() {
		yesterdayPrefix = "<"
	}
	tomorrowSuffix := ""
	if t.IsTomorrow() {
		tomorrowSuffix = ">"
	}
	hour, amPmSuffix := func() (int, string) {
		if t.format.Use24HourClock {
			return t.hour, ""
		}
		if t.hour == 12 {
			return 12, "pm"
		}
		if t.hour > 12 {
			return t.hour - 12, "pm"
		}
		if t.hour == 0 {
			return 12, "am"
		}
		return t.hour, "am"
	}()
	return fmt.Sprintf("%s%d:%02d%s%s", yesterdayPrefix, hour, t.minute, amPmSuffix, tomorrowSuffix)
}

func (t *time) ToStringWithFormat(f TimeFormat) string {
	c := *t
	c.format = f
	return c.ToString()
}

func (t *time) Format() TimeFormat {
	return t.format
}
