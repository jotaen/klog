package klog

import (
	"cloud.google.com/go/civil"
	"errors"
	"fmt"
	"regexp"
	"strconv"
	gotime "time"
)

type Time interface {
	Hour() int
	Minute() int
	MidnightOffset() Duration
	IsYesterday() bool
	IsTomorrow() bool
	IsToday() bool
	IsEqualTo(Time) bool
	IsAfterOrEqual(Time) bool
	Add(Duration) (Time, error)
	ToString() string
}

type time struct {
	hour          int
	minute        int
	dayShift      int
	is24HourClock bool
}

func newTime(hour int, minute int, dayShift int, is24HourClock bool) (Time, error) {
	ct := civil.Time{Hour: hour, Minute: minute}
	if !ct.IsValid() {
		return nil, errors.New("INVALID_TIME")
	}
	return &time{
		hour:          ct.Hour,
		minute:        ct.Minute,
		dayShift:      dayShift,
		is24HourClock: is24HourClock,
	}, nil
}

func NewTime(hour int, minute int) (Time, error) {
	return newTime(hour, minute, 0, true)
}

func NewTimeYesterday(hour int, minute int) (Time, error) {
	return newTime(hour, minute, -1, true)
}

func NewTimeTomorrow(hour int, minute int) (Time, error) {
	return newTime(hour, minute, +1, true)
}

var timePattern = regexp.MustCompile(`^(<)?(\d{1,2}):(\d{2})(am|pm)?(>)?$`)

func NewTimeFromString(hhmm string) (Time, error) {
	match := timePattern.FindStringSubmatch(hhmm)
	if len(match) != 6 || (match[1] == "<" && match[5] == ">") {
		return nil, errors.New("MALFORMED_TIME")
	}
	hour, _ := strconv.Atoi(match[2])
	minute, _ := strconv.Atoi(match[3])
	is24HourClock := true
	if match[4] == "am" || match[4] == "pm" {
		if hour < 1 || hour > 12 {
			return nil, errors.New("INVALID_TIME")
		}
		is24HourClock = false
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
	return newTime(hour, minute, dayShift, is24HourClock)
}

func NewTimeFromTime(t gotime.Time) Time {
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

func (t *time) Add(d Duration) (Time, error) {
	ONE_DAY := 24 * 60
	mins := t.MidnightOffset().Plus(d).InMinutes()
	if mins > 2*ONE_DAY || mins < ONE_DAY*-1 {
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
	result := &time{
		hour:          mins / 60,
		minute:        mins % 60,
		dayShift:      dayShift,
		is24HourClock: t.is24HourClock,
	}
	return result, nil
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
	hour := t.hour
	hour, am_pm := func() (int, string) {
		if t.is24HourClock {
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
	return fmt.Sprintf("%s%d:%02d%s%s", yesterdayPrefix, hour, t.minute, am_pm, tomorrowSuffix)
}
