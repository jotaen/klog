package parser

import (
	"github.com/jotaen/klog/klog/parser/engine"
	"github.com/jotaen/klog/klog/parser/txt"
)

func NewSerialParser() Parser {
	return serialParser
}

func NewParallelParser(numberOfWorkers int) Parser {
	return engine.ParallelParser[string, txt.Block, ParsedRecord, txt.Error]{
		Serialparser:    serialParser,
		NumberOfWorkers: numberOfWorkers,
	}
}

var serialParser = engine.SerialParser[string, txt.Block, ParsedRecord, txt.Error]{
	PreProcess: preProcess,
	ParseOne:   parse,
}
