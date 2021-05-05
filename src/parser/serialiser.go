package parser

import (
	"klog"
	"strings"
)

// SerialiseRecords serialises records into the canonical string representation.
func (h *Serialiser) SerialiseRecords(rs ...klog.Record) string {
	var text []string
	for _, r := range rs {
		text = append(text, h.serialiseRecord(r))
	}
	return strings.Join(text, "\n")
}

func (h *Serialiser) serialiseRecord(r klog.Record) string {
	text := ""
	text += h.Date(r.Date())
	if r.ShouldTotal().InMinutes() != 0 {
		text += " (" + h.ShouldTotal(r.ShouldTotal()) + ")"
	}
	text += "\n"
	if r.Summary() != "" {
		text += h.Summary(r.Summary()) + "\n"
	}
	for _, e := range r.Entries() {
		text += "    " // indentation
		text += (e.Unbox(
			func(r klog.Range) interface{} { return h.Range(r) },
			func(d klog.Duration) interface{} { return h.Duration(d, false) },
			func(o klog.OpenRange) interface{} { return h.OpenRange(o) },
		)).(string)
		if e.Summary() != "" {
			text += " " + h.Summary(e.Summary())
		}
		text += "\n"
	}
	return text
}

type Serialiser struct {
	Date        func(klog.Date) string
	ShouldTotal func(klog.Duration) string
	Summary     func(klog.Summary) string
	Range       func(klog.Range) string
	OpenRange   func(klog.OpenRange) string
	Duration    func(klog.Duration, bool) string
	Time        func(klog.Time) string
}

var DefaultSerialiser = Serialiser{
	Date:        func(d klog.Date) string { return d.ToString() },
	ShouldTotal: func(d klog.Duration) string { return d.ToString() },
	Summary:     func(s klog.Summary) string { return string(s) },
	Range:       func(r klog.Range) string { return r.ToString() },
	OpenRange:   func(or klog.OpenRange) string { return or.ToString() },
	Duration:    func(d klog.Duration, _ bool) string { return d.ToString() },
	Time:        func(t klog.Time) string { return t.ToString() },
}
