package parser

import (
	"github.com/jotaen/klog/klog"
	"github.com/jotaen/klog/klog/parser/engine"
	"github.com/jotaen/klog/klog/parser/txt"
)

// Parser parses a text into a list of Record datastructures. On success, it returns
// the parsed records. Otherwise, it returns all encountered parser errors.
type Parser interface {
	// Parse parses records from a string. It returns them along with the
	// respective blocks. Those two arrays have the same length.
	// Errors are reported via the last error array. In this case, the records
	// and blocks are nil. Note, one record can produce multiple errors,
	// so the length of the error array doesn’t say anything about the number
	// of records.
	Parse(string) ([]klog.Record, []txt.Block, []txt.Error)
}

// NewSerialParser returns a new parser, which processes the input text
// serially, i.e. one after the other.
func NewSerialParser() Parser {
	return serialParser
}

// NewParallelParser returns a new parser, which processes the input text
// in parallel. The parsing result is the same as with the serial parser.
func NewParallelParser(numberOfWorkers int) Parser {
	return engine.ParallelBatchParser[klog.Record]{
		SerialParser:    serialParser,
		NumberOfWorkers: numberOfWorkers,
	}
}

var serialParser = engine.SerialParser[klog.Record]{
	ParseOne: parse,
}
