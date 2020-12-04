package parser

import (
	"errors"
	"klog/entry"
)

const (
	MALFORMED_YAML = "The syntax of the document is not valid"
	DATE_MISSING   = "The date property must be set"
	INVALID_DATE   = "The date does not represent a valid day in the calendar"
	INVALID_TIME   = "The time"
	NEGATIVE_TIME  = "A time cannot be a negative value"
)

func parserError(code string) error {
	return errors.New(code)
}

func fromEntryError(err error) error {
	dict := map[string]string{
		entry.INVALID_DATE:  INVALID_DATE,
		entry.NEGATIVE_TIME: NEGATIVE_TIME,
	}
	return parserError(dict[err.(*entry.EntryError).Code])
}
