package datetime

import (
	"fmt"
)

type Date struct {
	Year  int
	Month int
	Day   int
}

func (d Date) ToString() string {
	return fmt.Sprintf("%04v-%02v-%02v", d.Year, d.Month, d.Day)
}

type Duration int64 // in minutes

type Time struct {
	Hour   int
	Minute int
}

func (t Time) ToString() string {
	return fmt.Sprintf("%02v:%02v", t.Hour, t.Minute)
}
