package service

import (
	. "github.com/jotaen/klog/src"
	"sort"
	gotime "time"
)

// Warning contains information for helping locate an issue.
type Warning struct {
	date   Date
	origin checker
}

// Date is the date of the record that the warning refers to.
func (w Warning) Date() Date {
	return w.date
}

// Warning is a short description of the problem.
func (w Warning) Warning() string {
	return w.origin.Message()
}

type checker interface {
	Warn(Record) Date
	Message() string
}

// CheckForWarnings checks records for potential user errors. It’s not meant as strict validation,
// but the main purpose is to help users spot accidental mistakes they might have made.
// The checks are mostly limited to record-level, because otherwise it would need to make
// assumptions on how records are organised within or across files.
func CheckForWarnings(reference gotime.Time, rs []Record) []Warning {
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
			d := c.Warn(r)
			if d != nil {
				ws = append(ws, Warning{
					date:   d,
					origin: c,
				})
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
func (c *unclosedOpenRangeChecker) Warn(record Record) Date {
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
		return record.Date()
	}
	return nil
}

func (c *unclosedOpenRangeChecker) Message() string {
	return "Unclosed open range"
}

type futureEntriesChecker struct {
	now         DateTime
	gracePeriod Duration
}

// Warn returns warnings if there are entries at future dates. It doesn’t
// return warnings if there are future records that don’t contain entries.
func (c *futureEntriesChecker) Warn(record Record) Date {
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
			countEntriesWithFutureTimes += Unbox[int](&e,
				func(r Range) int {
					if NewDateTime(record.Date(), r.Start()).IsAfterOrEqual(fuzzyNow) || NewDateTime(record.Date(), r.End()).IsAfterOrEqual(fuzzyNow) {
						return 1
					}
					return 0
				}, func(Duration) int {
					if record.Date().IsAfterOrEqual(c.now.Date.PlusDays(1)) {
						return 1
					}
					return 0
				}, func(or OpenRange) int {
					if NewDateTime(record.Date(), or.Start()).IsAfterOrEqual(fuzzyNow) {
						return 1
					}
					return 0
				})
		}
		if countEntriesWithFutureTimes == 0 {
			return nil
		}
	}
	return record.Date()
}

func (c *futureEntriesChecker) Message() string {
	return "Entry in the future"
}

type overlappingTimeRangesChecker struct{}

// Warn returns warnings if there are entries with overlapping time ranges.
// E.g. `8:00-9:00` and `8:30-9:30`.
func (c *overlappingTimeRangesChecker) Warn(record Record) Date {
	var orderedRanges []Range
	for _, e := range record.Entries() {
		Unbox(&e,
			func(r Range) any {
				orderedRanges = append(orderedRanges, r)
				return nil
			},
			func(Duration) any { return nil },
			func(or OpenRange) any {
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
			return record.Date()
		}
	}
	return nil
}

func (c *overlappingTimeRangesChecker) Message() string {
	return "Overlapping time ranges"
}

type moreThan24HoursChecker struct{}

// Warn returns warnings if there are records with a total time of more than 24h.
func (c *moreThan24HoursChecker) Warn(record Record) Date {
	if Total(record).InMinutes() > 24*60 {
		return record.Date()
	}
	return nil
}

func (c *moreThan24HoursChecker) Message() string {
	return "Total time exceeds 24 hours"
}
