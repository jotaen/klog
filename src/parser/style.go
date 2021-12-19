package parser

import . "github.com/jotaen/klog/src"

// Style describes the general styling and formatting preferences of a record.
type Style struct {
	lineEnding    string
	lineEndingSet bool

	indentation    string
	indentationSet bool

	spacingInRange    string // Example: `8:00 - 9:00` vs. `8:00-9:00`
	spacingInRangeSet bool

	dateFormat    DateFormat
	dateFormatSet bool

	timeFormat    TimeFormat
	timeFormatSet bool
}

func (s *Style) SetLineEnding(x string) {
	s.lineEnding = x
	s.lineEndingSet = true
}

func (s *Style) LineEnding() string {
	return s.lineEnding
}

func (s *Style) SetIndentation(x string) {
	s.indentation = x
	s.indentationSet = true
}

func (s *Style) Indentation() string {
	return s.indentation
}

func (s *Style) SetSpacingInRange(x string) {
	s.spacingInRange = x
	s.spacingInRangeSet = true
}

func (s *Style) SpacingInRange() string {
	return s.spacingInRange
}

func (s *Style) SetDateFormat(x DateFormat) {
	s.dateFormat = x
	s.dateFormatSet = true
}

func (s *Style) DateFormat() DateFormat {
	return s.dateFormat
}

func (s *Style) SetTimeFormat(x TimeFormat) {
	s.timeFormat = x
	s.timeFormatSet = true
}

func (s *Style) TimeFormat() TimeFormat {
	return s.timeFormat
}

// DefaultStyle returns the canonical style preferences as recommended
// by the file format specification.
func DefaultStyle() *Style {
	return &Style{
		lineEnding:     "\n",
		indentation:    "    ",
		spacingInRange: " ",
		dateFormat:     DateFormat{UseDashes: true},
		timeFormat:     TimeFormat{Use24HourClock: true},
	}
}

// Elect fills all unset fields of the `defaults` style with that value
// which was encountered most often in the parsed records. Fields of
// `defaults` that had been set explicitly take precedence.
func Elect(defaults Style, parsedRecords []ParsedRecord) *Style {
	lineEndingVotes := make(map[string]int)
	indentationVotes := make(map[string]int)
	spacingInRangeVotes := make(map[string]int)
	dateFormatVotes := make(map[DateFormat]int)
	timeFormatVotes := make(map[TimeFormat]int)
	for _, r := range parsedRecords {
		if r.Style.lineEndingSet {
			lineEndingVotes[r.Style.lineEnding] += 1
		}
		if r.Style.indentationSet {
			indentationVotes[r.Style.indentation] += 1
		}
		if r.Style.spacingInRangeSet {
			spacingInRangeVotes[r.Style.spacingInRange] += 1
		}
		if r.Style.dateFormatSet {
			dateFormatVotes[r.Style.dateFormat] += 1
		}
		if r.Style.timeFormatSet {
			timeFormatVotes[r.Style.timeFormat] += 1
		}
	}
	lineEndingMax := 0
	if !defaults.lineEndingSet {
		for x, v := range lineEndingVotes {
			if v > lineEndingMax {
				lineEndingMax = v
				defaults.SetLineEnding(x)
			}
		}
	}
	indentationMax := 0
	if !defaults.indentationSet {
		for x, v := range indentationVotes {
			if v > indentationMax {
				indentationMax = v
				defaults.SetIndentation(x)
			}
		}
	}
	spacingInRangeMax := 0
	if !defaults.spacingInRangeSet {
		for x, v := range spacingInRangeVotes {
			if v > spacingInRangeMax {
				spacingInRangeMax = v
				defaults.SetSpacingInRange(x)
			}
		}
	}
	dateFormatMax := 0
	if !defaults.dateFormatSet {
		for x, v := range dateFormatVotes {
			if v > dateFormatMax {
				dateFormatMax = v
				defaults.SetDateFormat(x)
			}
		}
	}
	timeFormatMax := 0
	if !defaults.timeFormatSet {
		for x, v := range timeFormatVotes {
			if v > timeFormatMax {
				timeFormatMax = v
				defaults.SetTimeFormat(x)
			}
		}
	}
	return &defaults
}
