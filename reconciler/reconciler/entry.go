package reconciler

import (
	"cloud.google.com/go/civil"
)

type Period struct {
	Start civil.Time
	End civil.Time
}

type Entry struct {
	Date civil.Date
	Summary string
	Times []int64
	Periods []Period
}

func (d Entry) TotalTime() (int64) {
	total := int64(0)
	for _, t := range d.Times {
		total += t
	}
	return total
}
