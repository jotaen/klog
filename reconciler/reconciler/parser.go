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

	entry := Entry{}

	date, err := civil.ParseDate(d.Date)
	if err == nil {
		entry.Date = date
	}

	if d.Summary != "" {
		entry.Summary = d.Summary
	}

	for _, h := range d.Hours {
		if h.Time != "" {
			t, _ := civil.ParseTime(h.Time + ":00")
			minutes := t.Minute + 60 * t.Hour
			entry.Times = append(entry.Times, Minutes(minutes))
		}
	}

	return entry, err
}
