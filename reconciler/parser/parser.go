package parser

import (
	"cloud.google.com/go/civil"
	"gopkg.in/yaml.v2"
	"main/entry"
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

func Parse(serialisedData string) (entry.Entry, error) {
	d, err := deserialise(serialisedData)
	if err != nil {
		return nil, err
	}

	date, err := civil.ParseDate(d.Date)
	if err != nil {
		return nil, err
	}
	e, _ := entry.Create(entry.Date{
		Year: date.Year,
		Month: date.Month,
		Day: date.Day,
	})

	e.SetSummary(d.Summary)

	for _, h := range d.Hours {
		if h.Time != "" {
			time, err := civil.ParseTime(h.Time + ":00")
			if err != nil {
				return nil, err
			}
			minutes := time.Minute + 60 * time.Hour
			e.AddTime(entry.Minutes(minutes))
		}
		if h.Start != "" && h.End != "" {
			start, _ := civil.ParseTime(h.Start + ":00")
			end, _ := civil.ParseTime(h.End + ":00")
			e.AddRange(
				entry.Time{ Hour: start.Hour, Minute: start.Minute },
				entry.Time{ Hour: end.Hour, Minute: end.Minute },
			)
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
