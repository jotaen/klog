package parser

import (
	. "klog/parser/parsing"
)

func ErrorInvalidDate(e Error) Error {
	return e.Set(
		"ErrorInvalidDate",
		"Invalid date",
		"Please make sure that the date format is either YYYY-MM-DD or YYYY/MM/DD, "+
			"and that its value represents a valid day in the calendar.",
	)
}

func ErrorIllegalIndentation(e Error) Error {
	return e.Set(
		"ErrorIllegalIndentation",
		"Unexpected indentation",
		"Please correct the indentation of this line.",
	)
}

func ErrorMalformedShouldTotal(e Error) Error {
	return e.Set(
		"ErrorMalformedShouldTotal",
		"Malformed property",
		"Please review the syntax of the should-total property. "+
			"Valid examples for it would be: (8h!) or (4h30m!) or (45m!)",
	)
}

func ErrorUnrecognisedProperty(e Error) Error {
	return e.Set(
		"ErrorUnrecognisedProperty",
		"Unrecognised property",
		"The highlighted property is not recognised. "+
			"The should-total property must be a time duration suffixed with an "+
			"exclamation mark, e.g. 5h15m! or 8h!",
	)
}

func ErrorMalformedPropertiesSyntax(e Error) Error {
	return e.Set(
		"ErrorMalformedPropertiesSyntax",
		"Malformed properties list",
		"Properties cannot be empty and must be "+
			"surrounded by parenthesis on both sides",
	)
}

func ErrorUnrecognisedTextInHeadline(e Error) Error {
	return e.Set(
		"ErrorUnrecognisedTextInHeadline",
		"Malformed headline",
		"The highlighted text in the headline is not recognised. "+
			"Please make sure to surround properties with parentheses, e.g.: (5h!) "+
			"You generally cannot put arbitrary text into the headline.",
	)
}

func ErrorMalformedSummary(e Error) Error {
	// this error cannot happen at the moment
	return e.Set(
		"ErrorMalformedSummary",
		"Invalid summary",
		"The summary text is not valid",
	)
}

func ErrorMalformedEntry(e Error) Error {
	return e.Set(
		"ErrorMalformedEntry",
		"Malformed entry",
		"Please review the syntax of the entry. "+
			"It must start with a duration or a time range. "+
			"Valid examples would be: 3h20m or 8:00-10:00 or 8:00-? "+
			"or <23:00-6:00 or 18:00-0:30>",
	)
}

func ErrorDuplicateOpenRange(e Error) Error {
	return e.Set(
		"ErrorDuplicateOpenRange",
		"Duplicate entry",
		"Please make sure that there is only "+
			"one open (unclosed) time range in this record.",
	)
}

func ErrorIllegalRange(e Error) Error {
	return e.Set(
		"ErrorIllegalRange",
		"Invalid date range",
		"Please make sure that both time values appear in chronological order. "+
			"If you want a time to be associated with an adjacent day you can use angle brackets "+
			"to shift the time by one day: <23:00-6:00 or 18:00-0:30>",
	)
}
