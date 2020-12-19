package datetime

import (
	"cloud.google.com/go/civil"
	"errors"
	"fmt"
	"regexp"
	gotime "time"
)

type Time interface {
	Hour() int
	Minute() int
	SinceMidnight() Duration
	ToString() string
}

type time struct {
	hour   int
	minute int
}

var timePattern = regexp.MustCompile(`^\s*\d{1,2}:\d{2}\s*$`)

func (t time) ToString() string {
	return fmt.Sprintf("%v:%02v", t.hour, t.minute)
}

func NewTime(hour int, minute int) (Time, error) {
	ct := civil.Time{
		Hour:       hour,
		Minute:     minute,
		Second:     0,
		Nanosecond: 0,
	}
	return ct2Time(ct)
}

func NewTimeFromString(hhmm string) (Time, error) {
	if !timePattern.MatchString(hhmm) {
		return nil, errors.New("MALFORMED_TIME")
	}
	ct, err := civil.ParseTime(hhmm + ":00")
	if err != nil {
		return nil, errors.New("INVALID_TIME")
	}
	return ct2Time(ct)
}

func CreateTimeFromTime(t gotime.Time) (Time, error) {
	return NewTime(t.Hour(), t.Minute())
}

func ct2Time(ct civil.Time) (Time, error) {
	if !ct.IsValid() {
		return nil, errors.New("INVALID_TIME")
	}
	return time{
		hour:   ct.Hour,
		minute: ct.Minute,
	}, nil
}

func (t time) Hour() int {
	return t.hour
}

func (t time) Minute() int {
	return t.minute
}

func (t time) SinceMidnight() Duration {
	return NewDuration(t.Hour(), t.Minute())
}
