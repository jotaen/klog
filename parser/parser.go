package parser

import (
	"errors"
	"gopkg.in/yaml.v2"
	"klog/datetime"
	"klog/workday"
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

func Parse(serialisedData string) (workday.WorkDay, []error) {
	errs := []error{}

	d, err := deserialise(serialisedData)
	if err != nil {
		errs = append(errs, errors.New("MALFORMED_YAML"))
		return nil, errs
	}

	date, err := datetime.CreateDateFromString(d.Date)
	if err != nil {
		errs = append(errs, err)
		return nil, errs
	}

	res := workday.Create(date)

	res.SetSummary(d.Summary)

	for _, h := range d.Hours {
		if h.Time != "" {
			time, err := datetime.CreateTimeFromString(h.Time)
			if err != nil {
				errs = append(errs, err)
			} else {
				minutes := time.Minute() + 60*time.Hour()
				res.AddTime(datetime.Duration(minutes))
			}
		}
		if h.Start != "" {
			start, _ := datetime.CreateTimeFromString(h.Start)
			if h.End != "" {
				end, _ := datetime.CreateTimeFromString(h.End)
				res.AddRange(start, end)
			} else {
				res.AddOpenRange(start)
			}
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
