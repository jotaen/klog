package parser

import (
	. "klog/parser/engine"
)

const (
	INVALID_VALUE        = "INVALID_VALUE"
	ILLEGAL_WHITESPACE   = "ILLEGAL_WHITESPACE"
	ILLEGAL_INDENTATION  = "ILLEGAL_INDENTATION"
	DUPLICATE_OPEN_RANGE = "DUPLICATE_OPEN_RANGE"
	ILLEGAL_RANGE        = "ILLEGAL_RANGE"
)

func ErrorMalformedDate(e Error) Error {
	return e.Set(INVALID_VALUE, "Date format is not okay")
}

func ErrorIllegalWhitespace(e Error) Error {
	return e.Set(ILLEGAL_WHITESPACE, "No whitespace allowed here")
}

func ErrorIllegalIndentation(e Error) Error {
	return e.Set(ILLEGAL_INDENTATION, "It is not allowed for a line with this indentation level to appear here")
}

func ErrorMalformedShouldTotal(e Error) Error {
	return e.Set(INVALID_VALUE, "Malformed should-total property")
}

func ErrorMalformedSummary(e Error) Error {
	return e.Set(INVALID_VALUE, "Please note that none of the lines of the summary is allowed to start with a whitespace character.")
}

func ErrorMalformedEntry(e Error) Error {
	return e.Set(INVALID_VALUE, "Malformed entry")
}

func ErrorDuplicateOpenRange(e Error) Error {
	return e.Set(DUPLICATE_OPEN_RANGE, "There cannot be two open time ranges in one record")
}

func ErrorIllegalRange(e Error) Error {
	return e.Set(ILLEGAL_RANGE, "The specified time range is not legal")
}
