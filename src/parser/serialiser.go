package parser

import (
	"klog/record"
)

func Serialise(rs []record.Record) string {
	text := ""
	for _, r := range rs {
		text += r.Date().ToString() + "\n"
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
			text += "\n"
		}
	}
	return text
}
