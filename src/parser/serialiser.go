package parser

import (
	"klog"
	"strings"
)

func SerialiseRecords(rs []src.Record, h FormattingHooks) string {
	var text []string
	if h == nil {
		h = defaultHooks{}
	}
	for _, r := range rs {
		text = append(text, SerialiseRecord(r, h))
	}
	return strings.Join(text, "\n")
}

func SerialiseRecord(r src.Record, h FormattingHooks) string {
	if h == nil {
		h = defaultHooks{}
	}
	text := ""
	text += h.PrintDate(r.Date())
	if r.ShouldTotal().InMinutes() != 0 {
		text += " (" + h.PrintShouldTotal(r.ShouldTotal(), "!") + ")"
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
		case src.Range:
			length = len(x.ToString())
			text += h.PrintRange(x)
		case src.Duration:
			length = len(x.ToString())
			text += h.PrintDuration(x)
		case src.OpenRange:
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
	PrintDate(src.Date) string
	PrintShouldTotal(src.Duration, string) string
	PrintSummary(src.Summary) string
	PrintRange(src.Range) string
	PrintOpenRange(src.OpenRange) string
	PrintDuration(src.Duration) string
}

type defaultHooks struct{}

func (h defaultHooks) PrintDate(d src.Date) string { return d.ToString() }
func (h defaultHooks) PrintShouldTotal(d src.Duration, symbol string) string {
	return d.ToString() + symbol
}
func (h defaultHooks) PrintSummary(s src.Summary) string      { return string(s) }
func (h defaultHooks) PrintRange(r src.Range) string          { return r.ToString() }
func (h defaultHooks) PrintOpenRange(or src.OpenRange) string { return or.ToString() }
func (h defaultHooks) PrintDuration(d src.Duration) string    { return d.ToString() }
