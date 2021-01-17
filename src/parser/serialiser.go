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
	maxLength := 0
	for _, e := range r.Entries() {
		length := len(e.ToString())
		if length > maxLength {
			maxLength = length
		}
	}
	for _, e := range r.Entries() {
		text += "    " // indentation
		length := 0
		switch x := e.Value().(type) {
		case Range:
			length = len(x.ToString())
			text += h.PrintRange(x)
		case Duration:
			length = len(x.ToString())
			text += h.PrintDuration(x)
		case OpenRange:
			length = len(x.ToString())
			text += h.PrintOpenRange(x)
		default:
			panic("Incomplete switch statement")
		}
		if e.Summary() != "" {
			text += strings.Repeat(" ", maxLength-length)
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
	PrintOpenRange(OpenRange) string
	PrintDuration(Duration) string
}

type defaultHooks struct{}

func (h defaultHooks) PrintDate(d Date) string            { return d.ToString() }
func (h defaultHooks) PrintShouldTotal(d Duration) string { return d.ToString() }
func (h defaultHooks) PrintSummary(s Summary) string      { return string(s) }
func (h defaultHooks) PrintRange(r Range) string          { return r.ToString() }
func (h defaultHooks) PrintOpenRange(or OpenRange) string { return or.ToString() }
func (h defaultHooks) PrintDuration(d Duration) string    { return d.ToString() }
