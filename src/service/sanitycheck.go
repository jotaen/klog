package service

import (
	. "github.com/jotaen/klog/src"
	"sort"
	gotime "time"
)

// Warning contains information for helping locate an issue.
type Warning struct {
	Date    Date
	Message string
}

type checker interface {
	Warn(Record) *Warning
}

// SanityCheck checks records for potential user errors. It’s not meant as strict validation,
// but the main purpose is to help users spot accidental mistakes they might have made.
func SanityCheck(reference gotime.Time, rs []Record) []Warning {
	now := NewDateTimeFromGo(reference)
	sortedRs := Sort(rs, false)
	var ws []Warning
	checkers := []checker{
		&unclosedOpenRangeChecker{today: now.Date},
		&futureEntriesChecker{now: now, gracePeriod: NewDuration(0, 31)},
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

// Warn returns warnings for all open ranges before yesterday, as these
// cannot be closed anymore via a shifted time. It also returns a warning
// if there is an open range yesterday, when there is a record today already.
func (c *unclosedOpenRangeChecker) Warn(record Record) *Warning {
	if record.Date().IsEqualTo(c.today) {
		// Open ranges at today’s date are always okay
		c.encounteredRecordAtToday = true
		return nil
	}
	if !c.encounteredRecordAtToday && c.today.PlusDays(-1).IsEqualTo(record.Date()) {
		// Open ranges at yesterday’s date are only okay if there is no entry today
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
	now         DateTime
	gracePeriod Duration
}

// Warn returns warnings if there are entries at future dates. It doesn’t
// return warnings if there are future records that don’t contain entries.
func (c *futureEntriesChecker) Warn(record Record) *Warning {
	if len(record.Entries()) == 0 {
		return nil
	}
	if c.now.Date.PlusDays(-2).IsAfterOrEqual(record.Date()) {
		return nil
	}
	if c.now.Date.PlusDays(-1).IsEqualTo(record.Date()) || c.now.Date.IsEqualTo(record.Date()) || c.now.Date.PlusDays(1).IsEqualTo(record.Date()) {
		countEntriesWithFutureTimes := 0
		fuzzyNow := func() DateTime {
			incTime, err := c.now.Time.Plus(c.gracePeriod)
			if err != nil {
				return c.now
			}
			return NewDateTime(c.now.Date, incTime)
		}()
		for _, e := range record.Entries() {
			countEntriesWithFutureTimes += e.Unbox(func(r Range) interface{} {
				if NewDateTime(record.Date(), r.Start()).IsAfterOrEqual(fuzzyNow) || NewDateTime(record.Date(), r.End()).IsAfterOrEqual(fuzzyNow) {
					return 1
				}
				return 0
			}, func(Duration) interface{} {
				if record.Date().IsAfterOrEqual(c.now.Date.PlusDays(1)) {
					return 1
				}
				return 0
			}, func(or OpenRange) interface{} {
				if NewDateTime(record.Date(), or.Start()).IsAfterOrEqual(fuzzyNow) {
					return 1
				}
				return 0
			}).(int)
		}
		if countEntriesWithFutureTimes == 0 {
			return nil
		}
	}
	return &Warning{
		Date:    record.Date(),
		Message: "Entry in the future",
	}
}

type overlappingTimeRangesChecker struct{}

// Warn returns warnings if there are entries with overlapping time ranges.
// E.g. `8:00-9:00` and `8:30-9:30`.
func (c *overlappingTimeRangesChecker) Warn(record Record) *Warning {
	var orderedRanges []Range
	for _, e := range record.Entries() {
		e.Unbox(
			func(r Range) interface{} {
				orderedRanges = append(orderedRanges, r)
				return nil
			},
			func(Duration) interface{} { return nil },
			func(or OpenRange) interface{} {
				// As best guess, assume open ranges to be closed at the end of the day.
				end, tErr := NewTime(23, 59)
				if tErr != nil {
					return nil
				}
				tr, rErr := NewRange(or.Start(), end)
				if rErr != nil {
					return nil
				}
				orderedRanges = append(orderedRanges, tr)
				return nil
			},
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

// Warn returns warnings if there are records with a total time of more than 24h.
func (c *moreThan24HoursChecker) Warn(record Record) *Warning {
	if Total(record).InMinutes() > 24*60 {
		return &Warning{
			Date:    record.Date(),
			Message: "Total time exceeds 24 hours",
		}
	}
	return nil
}
