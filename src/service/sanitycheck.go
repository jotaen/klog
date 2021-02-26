package service

import (
	. "klog"
	gotime "time"
)

type Warning struct {
	Date    Date
	Message string
}

type checker interface {
	Warn(Record) *Warning
}

// SanityCheck checks records for potential user errors.
func SanityCheck(reference gotime.Time, rs []Record) []Warning {
	today := NewDateFromTime(reference)
	sortedRs := Sort(rs, false)
	var ws []Warning
	checkers := []checker{
		&unclosedOpenRangeChecker{today: today},
		&futureEntriesChecker{today: today},
	}
	for _, r := range sortedRs {
		for _, c := range checkers {
			w := c.Warn(r)
			if w != nil {
				ws = append(ws, *w)
			}
		}
	}
	return ws
}

type unclosedOpenRangeChecker struct {
	today                    Date
	encounteredRecordAtToday bool
}

func (c *unclosedOpenRangeChecker) Warn(r Record) *Warning {
	if r.Date().IsEqualTo(c.today) {
		// Open ranges at today’s date are always okay
		c.encounteredRecordAtToday = true
		return nil
	}
	if !c.encounteredRecordAtToday && c.today.PlusDays(-1).IsEqualTo(r.Date()) {
		// Open ranges at yesterday’s date are only okay if there is no entry today today
		return nil
	}
	if r.OpenRange() != nil {
		// Any other case is most likely a mistake
		return &Warning{
			Date:    r.Date(),
			Message: "Unclosed open range",
		}
	}
	return nil
}

type futureEntriesChecker struct {
	today Date
}

func (c *futureEntriesChecker) Warn(r Record) *Warning {
	if r.Date().IsAfterOrEqual(c.today.PlusDays(1)) && len(r.Entries()) > 0 {
		return &Warning{
			Date:    r.Date(),
			Message: "Entry in future record",
		}
	}
	return nil
}
