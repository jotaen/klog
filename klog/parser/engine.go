package parser

import (
	"github.com/jotaen/klog/klog"
	"github.com/jotaen/klog/klog/parser/engine"
	"github.com/jotaen/klog/klog/parser/txt"
)

// Parser parses a text into a list of Record datastructures. On success, it returns
// the parsed records. Otherwise, it returns all encountered parser errors.
type Parser interface {
	Parse(string) ([]klog.Record, []txt.Block, []txt.Error)
}

func NewSerialParser() Parser {
	return serialParser
}

func NewParallelParser(numberOfWorkers int) Parser {
	if numberOfWorkers <= 0 {
		panic("ILLEGAL_WORKER_SIZE")
	}
	return engine.ParallelBatchParser[klog.Record]{
		NumberOfWorkers: numberOfWorkers,
		SerialParser:    serialParser,
	}
}

var serialParser = engine.SerialParser[klog.Record]{
	ParseOne: parse,
}
