package parser

import (
	"klog/record"
	"strings"
)

func SerialiseRecords(rs []record.Record) string {
	var text []string
	for _, r := range rs {
		text = append(text, SerialiseRecord(r))
	}
	return strings.Join(text, "\n")
}

func SerialiseRecord(r record.Record) string {
	text := ""
	text += r.Date().ToString()
	if r.ShouldTotal() != nil {
		text += " (" + r.ShouldTotal().ToString() + "!)"
	}
	text += "\n"
	if r.Summary() != "" {
		text += r.Summary() + "\n"
	}
	for _, e := range r.Entries() {
		text += "\t"
		switch x := e.Value().(type) {
		case record.Range:
			text += x.Start().ToString() + " - " + x.End().ToString()
			break
		case record.Duration:
			text += x.ToString()
			break
		case record.OpenRangeStart:
			text += x.ToString() + " -"
			break
		}
		if e.SummaryAsString() != "" {
			text += " " + e.SummaryAsString()
		}
		text += "\n"
	}
	return text
}
