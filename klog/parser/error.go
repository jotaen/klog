package parser

import "github.com/jotaen/klog/klog/parser/engine"

type HumanError struct {
	code    string
	title   string
	details string
}

func (e HumanError) New(t engine.Line, start int, length int) engine.Error {
	return engine.NewError(t, start, length, e.code, e.title, e.details)
}

func ErrorInvalidDate() HumanError {
	return HumanError{
		"ErrorInvalidDate",
		"Invalid date",
		"Please make sure that the date format is either YYYY-MM-DD or YYYY/MM/DD, " +
			"and that its value represents a valid day in the calendar.",
	}
}

func ErrorIllegalIndentation() HumanError {
	return HumanError{
		"ErrorIllegalIndentation",
		"Unexpected indentation",
		"Please correct the indentation of this line. Indentation must be 2-4 spaces or one tab. " +
			"You cannot mix different indentation styles within the same record.",
	}
}

func ErrorMalformedShouldTotal() HumanError {
	return HumanError{
		"ErrorMalformedShouldTotal",
		"Malformed should-total time",
		"Please review the syntax of the should-total time. " +
			"Valid examples for it would be: (8h!) or (4h30m!) or (45m!)",
	}
}

func ErrorUnrecognisedProperty() HumanError {
	return HumanError{
		"ErrorUnrecognisedProperty",
		"Unrecognised should-total value",
		"The highlighted value is not recognised. " +
			"The should-total must be a time duration suffixed with an " +
			"exclamation mark, e.g. 5h15m! or 8h!",
	}
}

func ErrorMalformedPropertiesSyntax() HumanError {
	return HumanError{
		"ErrorMalformedPropertiesSyntax",
		"Malformed should-total time",
		"The should-total cannot be empty and it must be " +
			"surrounded by parenthesis on both sides",
	}
}

func ErrorUnrecognisedTextInHeadline() HumanError {
	return HumanError{
		"ErrorUnrecognisedTextInHeadline",
		"Malformed headline",
		"The highlighted text in the headline is not recognised. " +
			"Please make sure to surround the should-total with parentheses, e.g.: (5h!) " +
			"You generally cannot put arbitrary text into the headline.",
	}
}

func ErrorMalformedSummary() HumanError {
	return HumanError{
		"ErrorMalformedSummary",
		"Malformed summary",
		"Summary lines cannot start with blank characters, such as non-breaking spaces.",
	}
}

func ErrorMalformedEntry() HumanError {
	return HumanError{
		"ErrorMalformedEntry",
		"Malformed entry",
		"Please review the syntax of the entry. " +
			"It must start with a duration or a time range. " +
			"Valid examples would be: 3h20m or 8:00-10:00 or 8:00-? " +
			"or <23:00-6:00 or 18:00-0:30>",
	}
}

func ErrorDuplicateOpenRange() HumanError {
	return HumanError{
		"ErrorDuplicateOpenRange",
		"Duplicate entry",
		"Please make sure that there is only " +
			"one open (unclosed) time range in this record.",
	}
}

func ErrorIllegalRange() HumanError {
	return HumanError{
		"ErrorIllegalRange",
		"Invalid date range",
		"Please make sure that both time values appear in chronological order. " +
			"If you want a time to be associated with an adjacent day you can use angle brackets " +
			"to shift the time by one day: <23:00-6:00 or 18:00-0:30>",
	}
}
