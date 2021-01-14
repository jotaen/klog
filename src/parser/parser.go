package parser

import (
	"fmt"
	"klog/datetime"
	"klog/record"
)

func Parse(serialisedData string) (record.Record, []ParserError) {
	errors := &errorCollection{}

	// Parse document
	d, err := parseYamlText(serialisedData)
	if err != nil {
		errors.add(parserError("MALFORMED_YAML", ""))
		return nil, errors.collection
	}

	// Parse date
	if d.Date == "" {
		errors.add(parserError("MISSING_REQUIRED_PROPERTY", "date"))
		return nil, errors.collection
	}
	date, err := datetime.NewDateFromString(d.Date)
	if err != nil {
		errors.add(fromError(err, fmt.Sprintf("date: %v", d.Date)))
		return nil, errors.collection
	}

	workDay := record.NewRecord(date)
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
			duration, err := datetime.NewDurationFromString(h.Time)
			if err != nil {
				errors.add(fromError(err, fmt.Sprintf("time: %v", h.Time)))
				continue
			}
			workDay.AddDuration(duration)
		}

		// Parse range
		if hasStart && hasEnd {
			// Start time
			start, err := datetime.NewTimeFromString(h.Start)
			if err != nil {
				errors.add(fromError(err, fmt.Sprintf("start: %v", h.Start)))
				continue
			}

			// End time
			end, err := datetime.NewTimeFromString(h.End)
			if err != nil {
				errors.add(fromError(err, fmt.Sprintf("end: %v", h.End)))
				continue
			}
			timeRange, err := datetime.NewTimeRange(start, end)
			if err != nil {
				errors.add(fromError(err, ""))
				continue
			}
			workDay.AddRange(timeRange)
		}

		// Parse open range
		if hasStart && !hasEnd {
			start, err := datetime.NewTimeFromString(h.Start)
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
