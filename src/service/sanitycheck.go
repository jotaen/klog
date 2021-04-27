package service

import (
	. "klog"
	"sort"
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
		&overlappingTimeRangesChecker{},
		&moreThan24HoursChecker{},
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

func (c *unclosedOpenRangeChecker) Warn(record Record) *Warning {
	if record.Date().IsEqualTo(c.today) {
		// Open ranges at today’s date are always okay
		c.encounteredRecordAtToday = true
		return nil
	}
	if !c.encounteredRecordAtToday && c.today.PlusDays(-1).IsEqualTo(record.Date()) {
		// Open ranges at yesterday’s date are only okay if there is no entry today today
		return nil
	}
	if record.OpenRange() != nil {
		// Any other case is most likely a mistake
		return &Warning{
			Date:    record.Date(),
			Message: "Unclosed open range",
		}
	}
	return nil
}

type futureEntriesChecker struct {
	today Date
}

func (c *futureEntriesChecker) Warn(record Record) *Warning {
	if record.Date().IsAfterOrEqual(c.today.PlusDays(1)) && len(record.Entries()) > 0 {
		return &Warning{
			Date:    record.Date(),
			Message: "Entry in future record",
		}
	}
	return nil
}

type overlappingTimeRangesChecker struct{}

func (c *overlappingTimeRangesChecker) Warn(record Record) *Warning {
	var orderedRanges []Range
	for _, e := range record.Entries() {
		e.Unbox(
			func(r Range) interface{} {
				orderedRanges = append(orderedRanges, r)
				return nil
			},
			func(Duration) interface{} { return nil },
			func(OpenRange) interface{} { return nil },
		)
	}
	sort.Slice(orderedRanges, func(i, j int) bool {
		return orderedRanges[j].Start().IsAfterOrEqual(orderedRanges[i].Start())
	})
	for i, curr := range orderedRanges {
		if i == 0 {
			continue
		}
		if curr.Start().IsEqualTo(curr.End()) {
			// Ignore point-in-time ranges
			continue
		}
		prev := orderedRanges[i-1]
		if !curr.Start().IsAfterOrEqual(prev.End()) {
			return &Warning{
				Date:    record.Date(),
				Message: "Overlapping time ranges",
			}
		}
	}
	return nil
}

type moreThan24HoursChecker struct{}

func (c *moreThan24HoursChecker) Warn(record Record) *Warning {
	if Total(record).InMinutes() > 24*60 {
		return &Warning{
			Date:    record.Date(),
			Message: "Total time exceeds 24 hours",
		}
	}
	return nil
}
