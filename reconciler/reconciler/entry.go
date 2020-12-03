package reconciler

import (
	"cloud.google.com/go/civil"
)

type Minutes int

type Period struct {
	Start civil.Time
	End civil.Time
}

type Entry struct {
	Date civil.Date
	Summary string
	Times []Minutes
	Periods []Period
}

func (d Entry) TotalTime() (Minutes) {
	total := Minutes(0)
	for _, t := range d.Times {
		total += t
	}
	return total
}
