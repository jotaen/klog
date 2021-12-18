package parser

import . "github.com/jotaen/klog/src"

// Style describes the general styling and formatting preferences of a record.
type Style struct {
	LineEnding    string
	lineEndingSet bool

	Indentation    string
	indentationSet bool

	SpacingInRange    string // Example: 8:00 - 9:00
	spacingInRangeSet bool

	DateFormat    DateFormat
	dateFormatSet bool

	TimeFormat    TimeFormat
	timeFormatSet bool
}

func (s *Style) SetLineEnding(x string) {
	s.LineEnding = x
	s.lineEndingSet = true
}

func (s *Style) SetIndentation(x string) {
	s.Indentation = x
	s.indentationSet = true
}

func (s *Style) SetSpacingInRange(x string) {
	s.SpacingInRange = x
	s.spacingInRangeSet = true
}

func (s *Style) SetDateFormat(x DateFormat) {
	s.DateFormat = x
	s.dateFormatSet = true
}

func (s *Style) SetTimeFormat(x TimeFormat) {
	s.TimeFormat = x
	s.timeFormatSet = true
}

// DefaultStyle returns the canonical style preferences as recommended
// by the file format specification.
func DefaultStyle() *Style {
	return &Style{
		LineEnding:     "\n",
		Indentation:    "    ",
		SpacingInRange: " ",
		DateFormat:     DateFormat{UseDashes: true},
		TimeFormat:     TimeFormat{Is24HourClock: true},
	}
}

func Elect(defaults Style, parsedRecords []ParsedRecord) *Style {
	lineEnding := make(map[string]int)
	lineEndingMax := 0
	indentation := make(map[string]int)
	indentationMax := 0
	spacingInRange := make(map[string]int)
	spacingInRangeMax := 0
	dateFormat := make(map[DateFormat]int)
	dateFormatMax := 0
	timeFormat := make(map[TimeFormat]int)
	timeFormatMax := 0
	for _, r := range parsedRecords {
		if r.Style.lineEndingSet {
			lineEnding[r.Style.LineEnding] += 1
		}
		if r.Style.indentationSet {
			indentation[r.Style.Indentation] += 1
		}
		if r.Style.spacingInRangeSet {
			spacingInRange[r.Style.SpacingInRange] += 1
		}
		if r.Style.dateFormatSet {
			dateFormat[r.Style.DateFormat] += 1
		}
		if r.Style.timeFormatSet {
			timeFormat[r.Style.TimeFormat] += 1
		}
	}
	if !defaults.lineEndingSet {
		for x, v := range lineEnding {
			if v > lineEndingMax {
				lineEndingMax = v
				defaults.SetLineEnding(x)
			}
		}
	}
	if !defaults.indentationSet {
		for x, v := range indentation {
			if v > indentationMax {
				indentationMax = v
				defaults.SetIndentation(x)
			}
		}
	}
	if !defaults.spacingInRangeSet {
		for x, v := range spacingInRange {
			if v > spacingInRangeMax {
				spacingInRangeMax = v
				defaults.SetSpacingInRange(x)
			}
		}
	}
	if !defaults.dateFormatSet {
		for x, v := range dateFormat {
			if v > dateFormatMax {
				dateFormatMax = v
				defaults.SetDateFormat(x)
			}
		}
	}
	if !defaults.timeFormatSet {
		for x, v := range timeFormat {
			if v > timeFormatMax {
				timeFormatMax = v
				defaults.SetTimeFormat(x)
			}
		}
	}
	return &defaults
}
