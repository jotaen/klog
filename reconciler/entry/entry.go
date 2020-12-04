package entry

import (
	"cloud.google.com/go/civil"
)

type Minutes int64

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

func (e Entry) Check() []EntryError {
	errs := []EntryError{}

	if !e.Date.IsValid() {
		errs = append(errs, EntryError{Code: INVALID_DATE})
	}

	for _, t := range e.Times {
		if t < 0 {
			errs = append(errs, EntryError{Code: NEGATIVE_TIME})
		}
	}

	if len(errs) == 0 {
		return nil
	}
	return errs
}

func (e Entry) TotalTime() Minutes {
	total := Minutes(0)
	for _, t := range e.Times {
		total += t
	}
	for _, t := range e.Ranges {
		start := t.Start.Minute + 60*t.Start.Hour
		end := t.End.Minute + 60*t.End.Hour
		total += Minutes(end - start)
	}
	return total
}
