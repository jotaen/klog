package parser

import (
	. "github.com/jotaen/klog/src"
	"strings"
)

// SerialiseRecords serialises records into the canonical string representation.
func (h *Serialiser) SerialiseRecords(rs ...Record) string {
	var text []string
	for _, r := range rs {
		text = append(text, h.serialiseRecord(r))
	}
	return strings.Join(text, "\n")
}

func (h *Serialiser) serialiseRecord(r Record) string {
	text := ""
	text += h.Date(r.Date())
	if r.ShouldTotal().InMinutes() != 0 {
		text += " (" + h.ShouldTotal(r.ShouldTotal()) + ")"
	}
	text += "\n"
	if !r.Summary().IsEmpty() {
		text += h.Summary(r.Summary()) + "\n"
	}
	for _, e := range r.Entries() {
		text += "    " // indentation
		text += (e.Unbox(
			func(r Range) interface{} { return h.Range(r) },
			func(d Duration) interface{} { return h.Duration(d) },
			func(o OpenRange) interface{} { return h.OpenRange(o) },
		)).(string)
		if !e.Summary().IsEmpty() {
			text += " " + h.Summary(e.Summary())
		}
		text += "\n"
	}
	return text
}

func serialiseSummary(s Summary) string {
	return strings.Join(s, "\n")
}

type Serialiser struct {
	Date           func(Date) string
	ShouldTotal    func(Duration) string
	Summary        func(Summary) string
	Range          func(Range) string
	OpenRange      func(OpenRange) string
	Duration       func(Duration) string
	SignedDuration func(Duration) string
	Time           func(Time) string
}

var PlainSerialiser = Serialiser{
	Date:           Date.ToString,
	ShouldTotal:    Duration.ToString,
	Summary:        serialiseSummary,
	Range:          Range.ToString,
	OpenRange:      OpenRange.ToString,
	Duration:       Duration.ToString,
	SignedDuration: Duration.ToStringWithSign,
	Time:           Time.ToString,
}
