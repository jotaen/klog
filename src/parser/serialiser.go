package parser

import (
	. "github.com/jotaen/klog/src"
	"strings"
)

// SerialiseRecords serialises records into the canonical string representation.
// (So it doesn’t and cannot restore the original formatting!)
func (h *Serialiser) SerialiseRecords(rs ...Record) string {
	var text []string
	for _, r := range rs {
		text = append(text, h.serialiseRecord(r))
	}
	return strings.Join(text, "\n")
}

var canonicalStyle = DefaultStyle()

func (h *Serialiser) serialiseRecord(r Record) string {
	text := ""
	text += h.Date(r.Date())
	if r.ShouldTotal().InMinutes() != 0 {
		text += " (" + h.ShouldTotal(r.ShouldTotal()) + ")"
	}
	text += canonicalStyle.LineEnding.Get()
	if r.Summary() != nil {
		text += h.Summary(SummaryText(r.Summary())) + canonicalStyle.LineEnding.Get()
	}
	for _, e := range r.Entries() {
		text += canonicalStyle.Indentation.Get()
		text += Unbox[string](&e,
			func(r Range) string { return h.Range(r) },
			func(d Duration) string { return h.Duration(d) },
			func(o OpenRange) string { return h.OpenRange(o) },
		)
		for i, l := range e.Summary().Lines() {
			if i == 0 && l != "" {
				text += " " // separator
			} else if i >= 1 {
				text += canonicalStyle.LineEnding.Get() + canonicalStyle.Indentation.Get() + canonicalStyle.Indentation.Get()
			}
			text += l
		}
		text += canonicalStyle.LineEnding.Get()
	}
	return text
}

type SummaryText []string

func (s SummaryText) ToString() string {
	return strings.Join(s, canonicalStyle.LineEnding.Get())
}

// Serialiser is used when the output should be modified, e.g. coloured.
type Serialiser struct {
	Date           func(Date) string
	ShouldTotal    func(Duration) string
	Summary        func(SummaryText) string
	Range          func(Range) string
	OpenRange      func(OpenRange) string
	Duration       func(Duration) string
	SignedDuration func(Duration) string
	Time           func(Time) string
}

// PlainSerialiser is used for unmodified (i.e. uncoloured) output.
var PlainSerialiser = Serialiser{
	Date:           Date.ToString,
	ShouldTotal:    Duration.ToString,
	Summary:        SummaryText.ToString,
	Range:          Range.ToString,
	OpenRange:      OpenRange.ToString,
	Duration:       Duration.ToString,
	SignedDuration: Duration.ToStringWithSign,
	Time:           Time.ToString,
}
