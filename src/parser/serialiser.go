package parser

import (
	"klog"
	"strings"
)

func SerialiseRecords(h *Serialiser, rs ...klog.Record) string {
	var text []string
	if h == nil {
		h = &defaultSerialiser
	}
	for _, r := range rs {
		text = append(text, serialiseRecord(h, r))
	}
	return strings.Join(text, "\n")
}

func serialiseRecord(h *Serialiser, r klog.Record) string {
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
}

var defaultSerialiser = Serialiser{
	Date:        func(d klog.Date) string { return d.ToString() },
	ShouldTotal: func(d klog.Duration) string { return d.ToString() },
	Summary:     func(s klog.Summary) string { return string(s) },
	Range:       func(r klog.Range) string { return r.ToString() },
	OpenRange:   func(or klog.OpenRange) string { return or.ToString() },
	Duration:    func(d klog.Duration, _ bool) string { return d.ToString() },
}
