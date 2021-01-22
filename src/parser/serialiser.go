package parser

import (
	"klog"
	"strings"
)

func SerialiseRecords(rs []klog.Record, h FormattingHooks) string {
	var text []string
	if h == nil {
		h = defaultHooks{}
	}
	for _, r := range rs {
		text = append(text, SerialiseRecord(r, h))
	}
	return strings.Join(text, "\n")
}

func SerialiseRecord(r klog.Record, h FormattingHooks) string {
	if h == nil {
		h = defaultHooks{}
	}
	text := ""
	text += h.PrintDate(r.Date())
	if r.ShouldTotal().InMinutes() != 0 {
		text += " (" + h.PrintShouldTotal(r.ShouldTotal()) + ")"
	}
	text += "\n"
	if r.Summary() != "" {
		text += h.PrintSummary(r.Summary()) + "\n"
	}
	for _, e := range r.Entries() {
		text += "    " // indentation
		text += (e.Unbox(
			func(r klog.Range) interface{} { return h.PrintRange(r) },
			func(d klog.Duration) interface{} { return h.PrintDuration(d) },
			func(o klog.OpenRange) interface{} { return h.PrintOpenRange(o) },
		)).(string)
		if e.Summary() != "" {
			text += " " + h.PrintSummary(e.Summary())
		}
		text += "\n"
	}
	return text
}

type FormattingHooks interface {
	PrintDate(klog.Date) string
	PrintShouldTotal(klog.Duration) string
	PrintSummary(klog.Summary) string
	PrintRange(klog.Range) string
	PrintOpenRange(klog.OpenRange) string
	PrintDuration(klog.Duration) string
}

type defaultHooks struct{}

func (h defaultHooks) PrintDate(d klog.Date) string            { return d.ToString() }
func (h defaultHooks) PrintShouldTotal(d klog.Duration) string { return d.ToString() }
func (h defaultHooks) PrintSummary(s klog.Summary) string      { return string(s) }
func (h defaultHooks) PrintRange(r klog.Range) string          { return r.ToString() }
func (h defaultHooks) PrintOpenRange(or klog.OpenRange) string { return or.ToString() }
func (h defaultHooks) PrintDuration(d klog.Duration) string    { return d.ToString() }
