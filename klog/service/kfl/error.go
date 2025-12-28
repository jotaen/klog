package kfl

import (
	"fmt"
	"math"
	"strings"

	tf "github.com/jotaen/klog/lib/text"
)

type ParseError interface {
	error
	Original() error
}

type parseError struct {
	err      error
	position int
	length   int
	query    string
}

func (e parseError) Error() string {
	errorLength := int(math.Max(float64(e.length), 1))
	relevantQueryFragment, newStart := tf.TextSubstrWithContext(e.query, e.position, errorLength, 10, 20)
	return fmt.Sprintf(
		"%s\n\n%s\n%s%s%s\n(Char %d in query.)",
		e.err,
		relevantQueryFragment,
		strings.Repeat("—", newStart),
		strings.Repeat("^", errorLength),
		strings.Repeat("—", len(relevantQueryFragment)-(newStart+errorLength)),
		e.position,
	)
}

func (e parseError) Original() error {
	return e.err
}
