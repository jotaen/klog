package record

import (
	"errors"
)

type Range interface {
	Start() Time
	End() Time
	Duration() Duration
	ToString() string
}

type OpenRange interface {
	Start() Time
	ToString() string
}

type timeRange struct {
	start Time
	end   Time
}

type openRange struct {
	start Time
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
	return tr.Start().ToString() + " - " + tr.End().ToString()
}

func (or *openRange) Start() Time {
	return or.start
}

func (or *openRange) ToString() string {
	return or.start.ToString() + " - ?"
}
