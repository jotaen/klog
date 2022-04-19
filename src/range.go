package klog

import (
	"errors"
)

// Range represents the period of time between two points of time.
type Range interface {
	Start() Time
	End() Time
	Duration() Duration

	// ToString serialises the range, e.g. `13:15 - 17:23`.
	ToString() string
}

// OpenRange represents a range that has not ended yet.
type OpenRange interface {
	Start() Time

	// ToString serialises the open range, e.g. `9:00 - ?`.
	ToString() string
}

func NewRange(start Time, end Time) (Range, error) {
	if !end.IsAfterOrEqual(start) {
		return nil, errors.New("ILLEGAL_RANGE")
	}
	return &timeRange{
		start: start,
		end:   end,
	}, nil
}

func NewOpenRange(start Time) OpenRange {
	return &openRange{start: start}
}

type timeRange struct {
	start Time
	end   Time
}

type openRange struct {
	start Time
}

func (tr *timeRange) Start() Time {
	return tr.start
}

func (tr *timeRange) End() Time {
	return tr.end
}

func (tr *timeRange) Duration() Duration {
	start := tr.Start().MidnightOffset().InMinutes()
	end := tr.End().MidnightOffset().InMinutes()
	return NewDuration(0, end-start)
}

func (tr *timeRange) ToString() string {
	return tr.Start().ToString() + " - " + tr.End().ToString() + " [" + tr.Duration().ToString() + "]"
}

func (or *openRange) Start() Time {
	return or.start
}

func (or *openRange) ToString() string {
	return or.start.ToString() + " - ?"
}
