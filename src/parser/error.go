package parser

import (
	. "klog/parser/engine"
)

const (
	INVALID_VALUE        = "INVALID_VALUE"
	ILLEGAL_SYNTAX       = "ILLEGAL_SYNTAX"
	UNRECOGNISED_TOKEN   = "UNRECOGNISED_TOKEN"
	ILLEGAL_WHITESPACE   = "ILLEGAL_WHITESPACE"
	ILLEGAL_INDENTATION  = "ILLEGAL_INDENTATION"
	DUPLICATE_OPEN_RANGE = "DUPLICATE_OPEN_RANGE"
	ILLEGAL_RANGE        = "ILLEGAL_RANGE"
)

func ErrorMalformedDate(e Error) Error {
	return e.Set(
		INVALID_VALUE,
		"Malformed date: please make sure that the date format is either YYYY-MM-DD or YYYY/MM/DD",
	)
}

func ErrorIllegalWhitespace(e Error) Error {
	return e.Set(
		ILLEGAL_WHITESPACE,
		"Illegal whitespace: please remove the highlighted whitespace characters",
	)
}

func ErrorIllegalIndentation(e Error, name string) Error {
	return e.Set(
		ILLEGAL_INDENTATION,
		"Unexpected indentation: please review the indentation of this line. "+
			"Expected an "+name+" to appear here.",
	)
}

func ErrorMalformedShouldTotal(e Error) Error {
	return e.Set(
		INVALID_VALUE,
		"Malformed property: please review the syntax of the should-total property. "+
			"Valid examples for it would be: (8h!) or (4h30m!) or (45m!)",
	)
}

func ErrorUnrecognisedProperty(e Error) Error {
	return e.Set(
		UNRECOGNISED_TOKEN,
		"Unrecognised property: the highlighted property is not recognised. "+
			"Please ensure that the should-total value must be suffixed with an "+
			"exclamation mark, e.g. (5h15m!)",
	)
}

func ErrorMalformedPropertiesSyntax(e Error) Error {
	return e.Set(
		ILLEGAL_SYNTAX,
		"Malformed properties list: please add a closing parenthesis so that the "+
			"properties are surrounded on both sides.",
	)
}

func ErrorUnrecognisedTextInHeadline(e Error) Error {
	return e.Set(
		ILLEGAL_SYNTAX,
		"Malformed headline: the highlighted text in the headline is not recognised. "+
			"Please make sure to surround properties with parentheses, e.g.: (5h!) "+
			"You generally cannot put arbitrary text into the headline.",
	)
}

func ErrorMalformedSummary(e Error) Error {
	return e.Set(
		INVALID_VALUE,
		"Illegal whitespace: please make sure that none of the lines of the summary "+
			"is allowed to start with a whitespace character.",
	)
}

func ErrorMalformedEntry(e Error) Error {
	return e.Set(
		INVALID_VALUE,
		"Malformed entry: please review the syntax of the entry. "+
			"It must start with a duration or a time range. "+
			"Valid examples would be: 3h20m or 8:00-10:00 or 8:00-? "+
			"or <23:00-6:00 or 18:00-0:30>",
	)
}

func ErrorDuplicateOpenRange(e Error) Error {
	return e.Set(
		DUPLICATE_OPEN_RANGE,
		"Invalid duplicate entry: please make sure that there is only "+
			"one open (unclosed) time range in this record.",
	)
}

func ErrorIllegalRange(e Error) Error {
	return e.Set(
		ILLEGAL_RANGE,
		"Invalid date range: please make sure that both time values appear in chronological order. "+
			"If you want a time to be associated with an adjacent day you can use angle brackets "+
			"to shift the time by one day: <23:00-6:00 or 18:00-0:30>",
	)
}
