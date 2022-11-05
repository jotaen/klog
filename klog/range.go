package klog

import (
	"errors"
	"strings"
)

// Range represents the period of time between two points of time.
type Range interface {
	Start() Time
	End() Time
	Duration() Duration

	// ToString serialises the range, e.g. `13:15 - 17:23`.
	ToString() string

	// Format returns the current formatting.
	Format() RangeFormat
}

// OpenRange represents a range that has not ended yet.
type OpenRange interface {
	Start() Time

	// ToString serialises the open range, e.g. `9:00 - ?`.
	ToString() string

	// Format returns the current formatting.
	Format() OpenRangeFormat
}

// RangeFormat contains the formatting options for a Range.
type RangeFormat struct {
	UseSpacesAroundDash bool
}

// DefaultRangeFormat returns the canonical time range format, as recommended by the spec.
func DefaultRangeFormat() RangeFormat {
	return RangeFormat{
		UseSpacesAroundDash: true,
	}
}

// OpenRangeFormat contains the formatting options for an OpenRange.
type OpenRangeFormat struct {
	UseSpacesAroundDash        bool
	AdditionalPlaceholderChars int
}

// DefaultOpenRangeFormat returns the canonical open range format, as recommended by the spec.
func DefaultOpenRangeFormat() OpenRangeFormat {
	return OpenRangeFormat{
		UseSpacesAroundDash:        DefaultRangeFormat().UseSpacesAroundDash,
		AdditionalPlaceholderChars: 0,
	}
}

func NewRange(start Time, end Time) (Range, error) {
	return NewRangeWithFormat(start, end, DefaultRangeFormat())
}

func NewRangeWithFormat(start Time, end Time, format RangeFormat) (Range, error) {
	if !end.IsAfterOrEqual(start) {
		return nil, errors.New("ILLEGAL_RANGE")
	}
	return &timeRange{
		start:  start,
		end:    end,
		format: format,
	}, nil
}

func NewOpenRange(start Time) OpenRange {
	return NewOpenRangeWithFormat(start, DefaultOpenRangeFormat())
}

func NewOpenRangeWithFormat(start Time, format OpenRangeFormat) OpenRange {
	return &openRange{start: start, format: format}
}

type timeRange struct {
	start  Time
	end    Time
	format RangeFormat
}

type openRange struct {
	start  Time
	format OpenRangeFormat
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
	space := " "
	if !tr.format.UseSpacesAroundDash {
		space = ""
	}
	return tr.Start().ToString() + space + "-" + space + tr.End().ToString()
}

func (tr *timeRange) Format() RangeFormat {
	return tr.format
}

func (or *openRange) Start() Time {
	return or.start
}

func (or *openRange) ToString() string {
	space := " "
	if !or.format.UseSpacesAroundDash {
		space = ""
	}
	return or.Start().ToString() + space + "-" + space + strings.Repeat("?", 1+or.format.AdditionalPlaceholderChars)
}

func (or *openRange) Format() OpenRangeFormat {
	return or.format
}
