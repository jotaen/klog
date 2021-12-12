package parser

import (
	. "github.com/jotaen/klog/src"
	. "github.com/jotaen/klog/src/parser/engine"
)

func ToRecords(pr []ParsedRecord) []Record {
	result := make([]Record, len(pr))
	for i, r := range pr {
		result[i] = r
	}
	return result
}

func ToBlocks(pr []ParsedRecord) []Block {
	result := make([]Block, len(pr))
	for i, r := range pr {
		result[i] = r.Block
	}
	return result
}
