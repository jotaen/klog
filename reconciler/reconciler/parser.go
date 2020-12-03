package reconciler

import (
	"cloud.google.com/go/civil"
	"gopkg.in/yaml.v2"
)

type data struct {
	Date    string
	Summary string
	Hours   []struct {
		Time  string
		Start string
		End   string
	}
}

func Parse(serialisedData string) (Entry, error) {
	d, err := deserialise(serialisedData)
	if err != nil {
		return Entry{}, err
	}

	e := Entry{}

	date, err := civil.ParseDate(d.Date)
	if err != nil {
		return Entry{}, err
	}
	e.Date = date

	if d.Summary != "" {
		e.Summary = d.Summary
	}

	for _, h := range d.Hours {
		if h.Time != "" {
			t, err := civil.ParseTime(h.Time + ":00")
			if err != nil {
				return Entry{}, err
			}
			minutes := t.Minute + 60*t.Hour
			e.Times = append(e.Times, Minutes(minutes))
		}
		if h.Start != "" && h.End != "" {
			start, _ := civil.ParseTime(h.Start + ":00")
			end, _ := civil.ParseTime(h.End + ":00")
			e.Ranges = append(e.Ranges, Range{Start: start, End: end})
		}
	}

	return e, nil
}

func deserialise(serialisedData string) (data, error) {
	d := data{}
	err := yaml.UnmarshalStrict([]byte(serialisedData), &d)
	if err != nil {
		return data{}, err
	}
	return d, nil
}
