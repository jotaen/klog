package parser

import . "klog/parser/engine"

const (
	INVALID_VALUE        = "INVALID_VALUE"
	ILLEGAL_WHITESPACE   = "ILLEGAL_WHITESPACE"
	ILLEGAL_INDENTATION  = "ILLEGAL_INDENTATION"
	DUPLICATE_OPEN_RANGE = "DUPLICATE_OPEN_RANGE"
	ILLEGAL_RANGE        = "ILLEGAL_RANGE"
)

func ErrorMalformedDate(e Error) Error {
	e.Code = INVALID_VALUE
	e.Message = "Date format is not okay"
	return e
}

func ErrorIllegalWhitespace(e Error) Error {
	e.Code = ILLEGAL_WHITESPACE
	e.Message = "No whitespace allowed here"
	return e
}

func ErrorIllegalIndentation(e Error) Error {
	e.Code = ILLEGAL_INDENTATION
	e.Message = "It is not allowed for a line with this indentation level to appear here"
	return e
}

func ErrorMalformedShouldTotal(e Error) Error {
	e.Code = INVALID_VALUE
	e.Message = "Malformed should-total property"
	return e
}

func ErrorMalformedSummary(e Error) Error {
	e.Code = INVALID_VALUE
	e.Message = "Please note that none of the lines of the summary is allowed to start with a whitespace character."
	return e
}

func ErrorMalformedEntry(e Error) Error {
	e.Code = INVALID_VALUE
	e.Message = "Malformed entry"
	return e
}

func ErrorDuplicateOpenRange(e Error) Error {
	e.Code = DUPLICATE_OPEN_RANGE
	e.Message = "There cannot be two open time ranges in one record"
	return e
}

func ErrorIllegalRange(e Error) Error {
	e.Code = ILLEGAL_RANGE
	e.Message = "The specified time range is not legal"
	return e
}
