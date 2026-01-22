package filter

import (
	"fmt"
	"math"
	"strings"

	tf "github.com/jotaen/klog/lib/terminalformat"
)

type ParseError interface {
	error
	Original() error
	Position() (int, int)
}

type parseError struct {
	err      error
	position int
	length   int
	query    string
}

func (e parseError) Error() string {
	errorLength := max(e.length, 1)
	relevantQueryFragment, newStart := tf.TextSubstrWithContext(e.query, e.position, errorLength, 10, 20)
	return fmt.Sprintf(
		"%s\n\n%s\n%s%s%s\nCursor positions %d-%d in query.",
		e.err,
		relevantQueryFragment,
		strings.Repeat("—", max(0, newStart)),
		strings.Repeat("^", max(0, errorLength)),
		strings.Repeat("—", max(0, len(relevantQueryFragment)-(newStart+errorLength))),
		e.position,
		e.position+errorLength,
	)
}

func (e parseError) Original() error {
	return e.err
}

func (e parseError) Position() (int, int) {
	return e.position, e.length
}

func max(x int, y int) int {
	return int(math.Max(float64(x), float64(y)))
}
