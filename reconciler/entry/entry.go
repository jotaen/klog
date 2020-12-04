package entry

import (
	"cloud.google.com/go/civil"
)

type Minutes int

type Range struct {
	Start civil.Time
	End   civil.Time
}

type Entry struct {
	Date    civil.Date
	Summary string
	Times   []Minutes
	Ranges  []Range
}

func (d Entry) TotalTime() Minutes {
	total := Minutes(0)
	for _, t := range d.Times {
		total += t
	}
	for _, t := range d.Ranges {
		start := t.Start.Minute + 60*t.Start.Hour
		end := t.End.Minute + 60*t.End.Hour
		total += Minutes(end - start)
	}
	return total
}
