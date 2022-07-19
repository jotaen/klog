package parser

import (
	. "github.com/jotaen/klog/src"
	"strings"
)

// Serialiser is used when the output should be modified, e.g. coloured.
type Serialiser interface {
	Date(Date) string
	ShouldTotal(Duration) string
	Summary(SummaryText) string
	Range(Range) string
	OpenRange(OpenRange) string
	Duration(Duration) string
	SignedDuration(Duration) string
	Time(Time) string
}

// SerialiseRecords serialises records into the canonical string representation.
// (So it doesn’t and cannot restore the original formatting!)
func SerialiseRecords(s Serialiser, rs ...Record) string {
	var text []string
	for _, r := range rs {
		text = append(text, serialiseRecord(s, r))
	}
	return strings.Join(text, "\n")
}

var canonicalStyle = DefaultStyle()

func serialiseRecord(s Serialiser, r Record) string {
	text := ""
	text += s.Date(r.Date())
	if r.ShouldTotal().InMinutes() != 0 {
		text += " (" + s.ShouldTotal(r.ShouldTotal()) + ")"
	}
	text += canonicalStyle.LineEnding.Get()
	if r.Summary() != nil {
		text += s.Summary(SummaryText(r.Summary())) + canonicalStyle.LineEnding.Get()
	}
	for _, e := range r.Entries() {
		text += canonicalStyle.Indentation.Get()
		text += Unbox[string](&e,
			func(r Range) string { return s.Range(r) },
			func(d Duration) string { return s.Duration(d) },
			func(o OpenRange) string { return s.OpenRange(o) },
		)
		for i, l := range e.Summary().Lines() {
			if i == 0 && l != "" {
				text += " " // separator
			} else if i >= 1 {
				text += canonicalStyle.LineEnding.Get() + canonicalStyle.Indentation.Get() + canonicalStyle.Indentation.Get()
			}
			text += s.Summary([]string{l})
		}
		text += canonicalStyle.LineEnding.Get()
	}
	return text
}

type SummaryText []string

func (s SummaryText) ToString() string {
	return strings.Join(s, canonicalStyle.LineEnding.Get())
}
