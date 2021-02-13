package service

import (
	. "klog"
	gotime "time"
)

type Warning struct {
	Date    Date
	Message string
}

func SanityCheck(reference gotime.Time, rs []Record) []Warning {
	today := NewDateFromTime(reference)
	var ws []Warning
	for _, r := range rs {
		if r.OpenRange() != nil && today.PlusDays(-2).IsAfterOrEqual(r.Date()) {
			ws = append(ws, Warning{
				Date:    r.Date(),
				Message: "Unclosed open range",
			})
		}
	}
	return ws
}
