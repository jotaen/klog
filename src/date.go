package klog

import (
	"cloud.google.com/go/civil"
	"errors"
	"fmt"
	"math"
	"regexp"
	"strings"
	gotime "time"
)

// Date represents a day in the gregorian calendar.
type Date interface {
	Year() int
	Month() int
	Day() int
	Weekday() int
	Quarter() int
	WeekNumber() int
	IsEqualTo(Date) bool
	IsAfterOrEqual(Date) bool
	ToString() string
	PlusDays(int) Date
}

type date struct {
	year             int
	month            int
	day              int
	formatWithDashes bool
}

var datePattern = regexp.MustCompile(`^(\d{4})[-/](\d{2})[-/](\d{2})$`)

func NewDate(year int, month int, day int) (Date, error) {
	cd := civil.Date{
		Year:  year,
		Month: gotime.Month(month),
		Day:   day,
	}
	return cd2Date(cd, true)
}

func NewDateFromString(yyyymmdd string) (Date, error) {
	match := datePattern.FindStringSubmatch(yyyymmdd)
	if len(match) != 4 || match[1] == "0" || match[2] == "0" || match[3] == "0" {
		return nil, errors.New("MALFORMED_DATE")
	}
	if c := strings.Count(yyyymmdd, "-"); c == 1 { // `-` and `/` mixed
		return nil, errors.New("MALFORMED_DATE")
	}
	cd, err := civil.ParseDate(match[1] + "-" + match[2] + "-" + match[3])
	if err != nil || !cd.IsValid() {
		return nil, errors.New("UNREPRESENTABLE_DATE")
	}
	return cd2Date(cd, strings.Contains(yyyymmdd, "-"))
}

func NewDateFromTime(t gotime.Time) Date {
	d, err := NewDate(t.Year(), int(t.Month()), t.Day())
	if err != nil {
		// This can/should never occur
		panic("ILLEGAL_DATE")
	}
	return d
}

func cd2Date(cd civil.Date, formatWithDashes bool) (Date, error) {
	if !cd.IsValid() {
		return nil, errors.New("UNREPRESENTABLE_DATE")
	}
	return &date{
		year:             cd.Year,
		month:            int(cd.Month),
		day:              cd.Day,
		formatWithDashes: formatWithDashes,
	}, nil
}

func (d *date) ToString() string {
	separator := "-"
	if !d.formatWithDashes {
		separator = "/"
	}
	return fmt.Sprintf("%04d%s%02d%s%02d", d.year, separator, d.month, separator, d.day)
}

func (d *date) Year() int {
	return d.year
}

func (d *date) Month() int {
	return d.month
}

func (d *date) Day() int {
	return d.day
}

func (d *date) Weekday() int {
	x := int(civil.Date{
		Year:  d.year,
		Month: gotime.Month(d.month),
		Day:   d.day,
	}.In(gotime.UTC).Weekday())
	if x == 0 {
		return 7
	}
	return x
}

func (d *date) Quarter() int {
	quarter := math.Ceil(float64(d.Month()) / 3)
	return int(quarter)
}

func (d *date) WeekNumber() int {
	_, week := civil.Date{
		Year:  d.year,
		Month: gotime.Month(d.month),
		Day:   d.day,
	}.In(gotime.UTC).ISOWeek()
	return week
}

func (d *date) IsEqualTo(otherDate Date) bool {
	return d.Year() == otherDate.Year() && d.Month() == otherDate.Month() && d.Day() == otherDate.Day()
}

func (d *date) IsAfterOrEqual(otherDate Date) bool {
	if d.Year() != otherDate.Year() {
		return d.Year() >= otherDate.Year()
	}
	if d.Month() != otherDate.Month() {
		return d.Month() >= otherDate.Month()
	}
	return d.Day() >= otherDate.Day()
}

func (d *date) PlusDays(dayIncrement int) Date {
	cd := civil.Date{
		Year:  d.year,
		Month: gotime.Month(d.month),
		Day:   d.day,
	}.AddDays(dayIncrement)
	newDate, err := cd2Date(cd, true)
	if err != nil {
		panic(err)
	}
	return newDate
}
