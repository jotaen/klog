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

//func NewParallelParser(numberOfWorkers int) Parser {
//	return engine.ParallelBatchParser[string, txt.Block, klog.Record, txt.Error]{
//		SerialParser:    serialParser,
//		NumberOfWorkers: numberOfWorkers,
//	}
//}

var serialParser = engine.SerialParser[klog.Record]{
	ParseOne: parse,
}
