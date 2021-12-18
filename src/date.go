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
	// Year returns the year as number, e.g. `2004`.
	Year() int

	// Month returns the month as number, e.g. `3` for March.
	Month() int

	// Day returns the day as number, e.g. `21`.
	Day() int

	// Weekday returns the day of the week, starting from Monday = 1.
	Weekday() int

	// Quarter returns the quarter that the date is in, e.g. `2` for `2010-04-15`.
	Quarter() int

	// WeekNumber returns the number of the week in the calendar year.
	WeekNumber() int

	// IsEqualTo checks whether two dates are the same.
	IsEqualTo(Date) bool

	// IsAfterOrEqual checks whether the given date occurs afterwards or at the same date.
	IsAfterOrEqual(Date) bool

	// PlusDays adds a number of days to the date. It doesn’t modify
	// the original object.
	PlusDays(int) Date

	// ToString serialises the date, e.g. `2017-04-23`.
	ToString() string

	ToStringWithFormat(DateFormat) string

	Format() DateFormat
}

type DateFormat struct {
	UseDashes bool
}

type date struct {
	year   int
	month  int
	day    int
	format DateFormat
}

var datePattern = regexp.MustCompile(`^(\d{4})[-/](\d{2})[-/](\d{2})$`)

func NewDate(year int, month int, day int) (Date, error) {
	cd := civil.Date{
		Year:  year,
		Month: gotime.Month(month),
		Day:   day,
	}
	return civil2Date(cd, DateFormat{UseDashes: true})
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
	return civil2Date(cd, DateFormat{UseDashes: strings.Contains(yyyymmdd, "-")})
}

func NewDateFromTime(t gotime.Time) Date {
	d, err := NewDate(t.Year(), int(t.Month()), t.Day())
	if err != nil {
		// This can/should never occur
		panic("ILLEGAL_DATE")
	}
	return d
}

func civil2Date(cd civil.Date, format DateFormat) (Date, error) {
	if !cd.IsValid() {
		return nil, errors.New("UNREPRESENTABLE_DATE")
	}
	if cd.Year > 9999 {
		// A year greater than 9999 cannot be serialised according to YYYY-MM-DD.
		return nil, errors.New("UNREPRESENTABLE_DATE")
	}
	return &date{
		year:   cd.Year,
		month:  int(cd.Month),
		day:    cd.Day,
		format: format,
	}, nil
}

func date2Civil(d *date) civil.Date {
	return civil.Date{
		Year:  d.year,
		Month: gotime.Month(d.month),
		Day:   d.day,
	}
}

func (d *date) ToString() string {
	separator := "-"
	if !d.format.UseDashes {
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
	x := int(date2Civil(d).In(gotime.UTC).Weekday())
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
	_, week := date2Civil(d).In(gotime.UTC).ISOWeek()
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
	cd := date2Civil(d).AddDays(dayIncrement)
	newDate, err := civil2Date(cd, d.format)
	if err != nil {
		panic(err)
	}
	return newDate
}

func (d *date) ToStringWithFormat(f DateFormat) string {
	nDate := *d
	nDate.format = f
	return nDate.ToString()
}

func (d *date) Format() DateFormat {
	return d.format
}
