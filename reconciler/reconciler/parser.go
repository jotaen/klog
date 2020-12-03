package reconciler

import (
	"cloud.google.com/go/civil"
	"gopkg.in/yaml.v2"
)

type data struct {
	Date string
	Summary string
	Hours []struct {
		Time string
		Start string
		End string
	}
}

func Parse(serialisedData string) (Entry, error) {
	d := data{}
	err := yaml.Unmarshal([]byte(serialisedData), &d)
	e := Entry{}

	date, err := civil.ParseDate(d.Date)
	if err != nil {
		return e, err
	}
	e.Date = date

	if d.Summary != "" {
		e.Summary = d.Summary
	}

	for _, h := range d.Hours {
		if h.Time != "" {
			t, _ := civil.ParseTime(h.Time + ":00")
			minutes := t.Minute + 60 * t.Hour
			e.Times = append(e.Times, Minutes(minutes))
		}
		if h.Start != "" && h.End != "" {
			start, _ := civil.ParseTime(h.Start + ":00")
			end, _ := civil.ParseTime(h.End + ":00")
			e.Ranges = append(e.Ranges, Range{ Start: start, End: end })
		}
	}

	return e, nil
}
