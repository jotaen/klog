package record

import (
	"cloud.google.com/go/civil"
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"strings"
	gotime "time"
)

type Time interface {
	Hour() int
	Minute() int
	MidnightOffset() Duration
	IsYesterday() bool
	IsTomorrow() bool
	IsToday() bool
	IsAfterOrEqual(Time) bool
	ToString() string
}

type time struct {
	hour     int
	minute   int
	dayShift int8
}

var timePattern = regexp.MustCompile(`^<?\d{1,2}:\d{2}>?$`)

func (t *time) ToString() string {
	prefix := ""
	if t.IsYesterday() {
		prefix = "<"
	}
	suffix := ""
	if t.IsTomorrow() {
		suffix = ">"
	}
	return fmt.Sprintf("%s%v:%02v%s", prefix, t.hour, t.minute, suffix)
}

func newTime(hour int, minute int, dayShift int8) (Time, error) {
	ct := civil.Time{Hour: hour, Minute: minute}
	if !ct.IsValid() {
		return nil, errors.New("INVALID_TIME")
	}
	return &time{
		hour:     ct.Hour,
		minute:   ct.Minute,
		dayShift: dayShift,
	}, nil
}

func NewTime(hour int, minute int) (Time, error) {
	return newTime(hour, minute, 0)
}

func NewTimeYesterday(hour int, minute int) (Time, error) {
	return newTime(hour, minute, -1)
}

func NewTimeTomorrow(hour int, minute int) (Time, error) {
	return newTime(hour, minute, +1)
}

func NewTimeFromString(hhmm string) (Time, error) {
	if !timePattern.MatchString(hhmm) || (strings.HasPrefix(hhmm, "<") && strings.HasSuffix(hhmm, ">")) {
		return nil, errors.New("MALFORMED_TIME")
	}
	dayShift := int8(0)
	if strings.HasPrefix(hhmm, "<") {
		dayShift = -1
		hhmm = strings.TrimPrefix(hhmm, "<")
	} else if strings.HasSuffix(hhmm, ">") {
		dayShift = +1
		hhmm = strings.TrimSuffix(hhmm, ">")
	}
	hhmm = strings.TrimSpace(hhmm)
	parts := strings.Split(hhmm, ":")
	hour, _ := strconv.Atoi(parts[0])
	minute, _ := strconv.Atoi(parts[1])
	return newTime(hour, minute, dayShift)
}

func CreateTimeFromTime(t gotime.Time) (Time, error) {
	return NewTime(t.Hour(), t.Minute())
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

func (t *time) IsAfterOrEqual(otherTime Time) bool {
	first := t.MidnightOffset()
	second := otherTime.MidnightOffset()
	return first.InMinutes() >= second.InMinutes()
}
