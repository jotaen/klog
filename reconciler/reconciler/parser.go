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

	date, err := civil.ParseDate(d.Date)
	entry := Entry{
		Date: date,
		Summary: d.Summary,
	}

	return entry, err
}
