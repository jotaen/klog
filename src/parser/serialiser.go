package parser

import (
	. "klog/record"
	"strings"
)

func SerialiseRecords(rs []Record, h FormattingHooks) string {
	var text []string
	if h == nil {
		h = defaultHooks{}
	}
	for _, r := range rs {
		text = append(text, SerialiseRecord(r, h))
	}
	return strings.Join(text, "\n")
}

func SerialiseRecord(r Record, h FormattingHooks) string {
	if h == nil {
		h = defaultHooks{}
	}
	text := ""
	text += h.PrintDate(r.Date())
	if r.ShouldTotal() != nil {
		text += " (" + h.PrintShouldTotal(r.ShouldTotal()) + "!)"
	}
	text += "\n"
	if r.Summary() != "" {
		text += h.PrintSummary(r.Summary()) + "\n"
	}
	for _, e := range r.Entries() {
		text += "\t"
		switch x := e.Value().(type) {
		case Range:
			text += h.PrintRange(x)
			break
		case Duration:
			text += h.PrintDuration(x)
			break
		case OpenRangeStart:
			text += x.ToString() + " -"
			break
		}
		if e.Summary() != "" {
			text += " " + h.PrintSummary(e.Summary())
		}
		text += "\n"
	}
	return text
}

type FormattingHooks interface {
	PrintDate(Date) string
	PrintShouldTotal(Duration) string
	PrintSummary(Summary) string
	PrintRange(Range) string
	PrintDuration(Duration) string
}

type defaultHooks struct{}

func (h defaultHooks) PrintDate(d Date) string            { return d.ToString() }
func (h defaultHooks) PrintShouldTotal(d Duration) string { return d.ToString() }
func (h defaultHooks) PrintSummary(s Summary) string      { return string(s) }
func (h defaultHooks) PrintRange(r Range) string          { return r.ToString() }
func (h defaultHooks) PrintDuration(d Duration) string    { return d.ToString() }
