package parser

import klog "github.com/jotaen/klog/src"

// Style describes the general styling and formatting preferences of a record.
type Style struct {
	LineEnding     string
	Indentation    string
	SpacingInRange string // Example: 8:00 - 9:00
	DateFormat     klog.DateFormat
	TimeFormat     klog.TimeFormat
}

// DefaultStyle returns the canonical style preferences as recommended
// by the file format specification.
func DefaultStyle() Style {
	return Style{
		LineEnding:     "\n",
		Indentation:    "    ",
		SpacingInRange: " ",
		DateFormat:     klog.DateFormat{UseDashes: true},
		TimeFormat:     klog.TimeFormat{Is24HourClock: true},
	}
}
