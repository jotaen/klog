package datetime

import (
	"cloud.google.com/go/civil"
	"errors"
	"fmt"
	gotime "time"
)

type Date interface {
	Year() int
	Month() int
	Day() int
	ToString() string
}

type date struct {
	year  int
	month int
	day   int
}

func CreateDate(year int, month int, day int) (Date, error) {
	cd := civil.Date{
		Year:  year,
		Month: gotime.Month(month),
		Day:   day,
	}
	return cd2Date(cd)
}

func CreateDateFromString(yyyymmdd string) (Date, error) {
	cd, err := civil.ParseDate(yyyymmdd)
	if err != nil {
		return nil, errors.New(INVALID_DATE)
	}
	return cd2Date(cd)
}

func cd2Date(cd civil.Date) (Date, error) {
	if !cd.IsValid() {
		return nil, errors.New(INVALID_DATE)
	}
	return date{
		year:  cd.Year,
		month: int(cd.Month),
		day:   cd.Day,
	}, nil
}

func (d date) ToString() string {
	return fmt.Sprintf("%04v-%02v-%02v", d.year, d.month, d.day)
}

func (d date) Year() int {
	return d.year
}

func (d date) Month() int {
	return d.month
}

func (d date) Day() int {
	return d.day
}
