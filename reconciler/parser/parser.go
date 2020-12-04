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

func Parse(serialisedData string) (entry.Entry, []error) {
	errs := []error{}

	d, err := deserialise(serialisedData)
	if err != nil {
		errs = append(errs, parserError(MALFORMED_YAML))
		return nil, errs
	}

	date, _ := civil.ParseDate(d.Date)
	res, err := entry.Create(entry.Date{
		Year: date.Year,
		Month: date.Month,
		Day: date.Day,
	})
	if res == nil {
		errs = append(errs, fromEntryError(err))
		return nil, errs
	}

	res.SetSummary(d.Summary)

	for _, h := range d.Hours {
		if h.Time != "" {
			time, err := civil.ParseTime(h.Time + ":00")
			if err != nil {
				errs = append(errs, parserError(INVALID_TIME))
			}
			minutes := time.Minute + 60 * time.Hour
			res.AddTime(entry.Minutes(minutes))
		}
		if h.Start != "" && h.End != "" {
			start, _ := civil.ParseTime(h.Start + ":00")
			end, _ := civil.ParseTime(h.End + ":00")
			res.AddRange(
				entry.Time{ Hour: start.Hour, Minute: start.Minute },
				entry.Time{ Hour: end.Hour, Minute: end.Minute },
			)
		}
	}

	if len(errs) != 0 {
		return nil, errs
	}
	return res, nil
}

func deserialise(serialisedData string) (data, error) {
	d := data{}
	err := yaml.UnmarshalStrict([]byte(serialisedData), &d)
	if err != nil {
		return data{}, err
	}
	return d, nil
}
