package filter

import (
	"math"
)

type ParseError interface {
	error
	Original() error
	Position() (int, int)
	Query() string
}

type parseError struct {
	err      error
	position int
	length   int
	query    string
}

func NewParseError() ParseError {
	return parseError{}
}

func (e parseError) Error() string {
	return "Illegal filter expression"
}

func (e parseError) Original() error {
	return e.err
}

func (e parseError) Query() string {
	return e.query
}

func (e parseError) Position() (int, int) {
	return e.position, e.length
}

func max(x int, y int) int {
	return int(math.Max(float64(x), float64(y)))
}
