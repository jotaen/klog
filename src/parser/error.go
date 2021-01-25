package parser

import (
	. "klog/parser/engine"
)

func ErrorInvalidDate(e Error) Error {
	return e.Set(
		"ErrorInvalidDate",
		"Invalid date: please make sure that the date format is either YYYY-MM-DD or YYYY/MM/DD, "+
			"and that its value represents a valid day in the calendar.",
	)
}

func ErrorIllegalIndentation(e Error) Error {
	return e.Set(
		"ErrorIllegalIndentation",
		"Unexpected indentation: please correct the indentation of this line.",
	)
}

func ErrorMalformedShouldTotal(e Error) Error {
	return e.Set(
		"ErrorMalformedShouldTotal",
		"Malformed property: please review the syntax of the should-total property. "+
			"Valid examples for it would be: (8h!) or (4h30m!) or (45m!)",
	)
}

func ErrorUnrecognisedProperty(e Error) Error {
	return e.Set(
		"ErrorUnrecognisedProperty",
		"Unrecognised property: the highlighted property is not recognised. "+
			"Please ensure that the should-total value must be suffixed with an "+
			"exclamation mark, e.g. (5h15m!)",
	)
}

func ErrorMalformedPropertiesSyntax(e Error) Error {
	return e.Set(
		"ErrorMalformedPropertiesSyntax",
		"Malformed properties list: properties cannot be empty and must be "+
			"surrounded by parenthesis on both sides",
	)
}

func ErrorUnrecognisedTextInHeadline(e Error) Error {
	return e.Set(
		"ErrorUnrecognisedTextInHeadline",
		"Malformed headline: the highlighted text in the headline is not recognised. "+
			"Please make sure to surround properties with parentheses, e.g.: (5h!) "+
			"You generally cannot put arbitrary text into the headline.",
	)
}

func ErrorMalformedSummary(e Error) Error {
	// this error cannot happen at the moment
	return e.Set(
		"ErrorMalformedSummary",
		"Invalid summary",
	)
}

func ErrorMalformedEntry(e Error) Error {
	return e.Set(
		"ErrorMalformedEntry",
		"Malformed entry: please review the syntax of the entry. "+
			"It must start with a duration or a time range. "+
			"Valid examples would be: 3h20m or 8:00-10:00 or 8:00-? "+
			"or <23:00-6:00 or 18:00-0:30>",
	)
}

func ErrorDuplicateOpenRange(e Error) Error {
	return e.Set(
		"ErrorDuplicateOpenRange",
		"Invalid duplicate entry: please make sure that there is only "+
			"one open (unclosed) time range in this record.",
	)
}

func ErrorIllegalRange(e Error) Error {
	return e.Set(
		"ErrorIllegalRange",
		"Invalid date range: please make sure that both time values appear in chronological order. "+
			"If you want a time to be associated with an adjacent day you can use angle brackets "+
			"to shift the time by one day: <23:00-6:00 or 18:00-0:30>",
	)
}
