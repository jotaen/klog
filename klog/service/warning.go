package service

import (
	"github.com/jotaen/klog/klog"
	"sort"
	gotime "time"
)

// Warning contains information for helping locate an issue.
type Warning struct {
	date   klog.Date
	origin checker
}

// Date is the date of the record that the warning refers to.
func (w Warning) Date() klog.Date {
	return w.date
}

// Warning is a short description of the problem.
func (w Warning) Warning() string {
	return w.origin.Message()
}

type checker interface {
	Warn(klog.Record) klog.Date
	Message() string
}

// CheckForWarnings checks records for potential logical issues in the data. For every
// issue encountered, it invokes the `onWarn` callback. Note: Warnings are not meant as
// strict validation, but the main purpose is to help users spot accidental mistakes users
// might have made. The checks are limited to record-level, because otherwise it would
// need to make assumptions on how records are organised within or across files.
func CheckForWarnings(onWarn func(Warning), reference gotime.Time, rs []klog.Record) {
	now := NewDateTimeFromGo(reference)
	sortedRs := Sort(rs, false)
	checkers := []checker{
		&unclosedOpenRangeChecker{today: now.Date},
		&futureEntriesChecker{now: now, gracePeriod: klog.NewDuration(0, 31)},
		&overlappingTimeRangesChecker{},
		&moreThan24HoursChecker{},
	}
	for _, r := range sortedRs {
		for _, c := range checkers {
			d := c.Warn(r)
			if d != nil {
				onWarn(Warning{
					date:   d,
					origin: c,
				})
			}
		}
	}
}

type unclosedOpenRangeChecker struct {
	today                    klog.Date
	encounteredRecordAtToday bool
}

// Warn returns warnings for all open ranges before yesterday, as these
// cannot be closed anymore via a shifted time. It also returns a warning
// if there is an open range yesterday, when there is a record today already.
func (c *unclosedOpenRangeChecker) Warn(record klog.Record) klog.Date {
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
	gracePeriod klog.Duration
}

// Warn returns warnings if there are entries at future dates. It doesn’t
// return warnings if there are future records that don’t contain entries.
func (c *futureEntriesChecker) Warn(record klog.Record) klog.Date {
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
			countEntriesWithFutureTimes += klog.Unbox[int](&e,
				func(r klog.Range) int {
					if NewDateTime(record.Date(), r.Start()).IsAfterOrEqual(fuzzyNow) || NewDateTime(record.Date(), r.End()).IsAfterOrEqual(fuzzyNow) {
						return 1
					}
					return 0
				}, func(klog.Duration) int {
					if record.Date().IsAfterOrEqual(c.now.Date.PlusDays(1)) {
						return 1
					}
					return 0
				}, func(or klog.OpenRange) int {
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
func (c *overlappingTimeRangesChecker) Warn(record klog.Record) klog.Date {
	var orderedRanges []klog.Range
	for _, e := range record.Entries() {
		klog.Unbox(&e,
			func(r klog.Range) any {
				orderedRanges = append(orderedRanges, r)
				return nil
			},
			func(klog.Duration) any { return nil },
			func(or klog.OpenRange) any {
				// As best guess, assume open ranges to be closed at the end of the day.
				end, tErr := klog.NewTime(23, 59)
				if tErr != nil {
					return nil
				}
				tr, rErr := klog.NewRange(or.Start(), end)
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
func (c *moreThan24HoursChecker) Warn(record klog.Record) klog.Date {
	if Total(record).InMinutes() > 24*60 {
		return record.Date()
	}
	return nil
}

func (c *moreThan24HoursChecker) Message() string {
	return "Total time exceeds 24 hours"
}
