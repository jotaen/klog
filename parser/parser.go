package parser

import (
	"fmt"
	"klog/datetime"
	"klog/workday"
	"strings"
)

func Parse(serialisedData string) (workday.WorkDay, []ParserError) {
	errors := &errorCollection{}

	// Parse document
	d, err := parseYamlText(serialisedData)
	if err != nil {
		errors.add(parserError("MALFORMED_YAML", ""))
		return nil, errors.collection
	}

	// Parse date
	date, err := datetime.CreateDateFromString(d.Date)
	if err != nil {
		errors.add(fromError(err, fmt.Sprintf("date: %v", d.Date)))
		return nil, errors.collection
	}

	workDay := workday.Create(date)
	workDay.SetSummary(d.Summary)

	// Parse hours
	for _, h := range d.Hours {
		hasTime := h.Time != ""
		hasStart := h.Start != ""
		hasEnd := h.End != ""
		if (hasTime && (hasStart || hasEnd)) || (hasEnd && !hasStart) {
			errors.add(parserError("MALFORMED_HOURS", "hours"))
			continue
		}

		// Parse time
		if hasTime {
			duration, err := datetime.CreateDurationFromString(h.Time)
			if err != nil {
				errors.add(fromError(err, fmt.Sprintf("time: %v", h.Time)))
				continue
			}
			workDay.AddDuration(duration)
		}

		// Parse range
		if hasStart && hasEnd {
			// Start time
			startTime := strings.Split(h.Start, " ")
			start, err := datetime.CreateTimeFromString(startTime[0])
			if err != nil {
				errors.add(fromError(err, fmt.Sprintf("start: %v", h.Start)))
				continue
			}
			isStartYesterday := false
			if len(startTime) == 2 && startTime[1] == "yesterday" {
				isStartYesterday = true
			}

			// End time
			endTime := strings.Split(h.End, " ")
			end, err := datetime.CreateTimeFromString(endTime[0])
			isEndTomorrow := false
			if len(endTime) == 2 && endTime[1] == "tomorrow" {
				isEndTomorrow = true
			}
			if err != nil {
				errors.add(fromError(err, fmt.Sprintf("end: %v", h.End)))
				continue
			}
			timeRange, err := datetime.CreateOverlappingTimeRange(start, isStartYesterday, end, isEndTomorrow)
			if err != nil {
				errors.add(fromError(err, ""))
				continue
			}
			workDay.AddRange(timeRange)
		}

		// Parse open range
		if hasStart && !hasEnd {
			start, err := datetime.CreateTimeFromString(h.Start)
			if err != nil {
				errors.add(fromError(err, fmt.Sprintf("start: %v", h.Start)))
				continue
			}
			workDay.StartOpenRange(start)
		}
	}

	if len(errors.collection) != 0 {
		return nil, errors.collection
	}
	return workDay, nil
}
