package parser

import (
	"github.com/jotaen/klog/klog/parser/txt"
)

type errData struct {
	id     string
	lineNr int
	pos    int
	len    int
}

func (e HumanError) toErrData(lineNr int, pos int, len int) errData {
	return errData{e.code, lineNr, pos, len}
}

func toErrData(e txt.Error) errData {
	return errData{e.Code(), e.LineNumber(), e.Position(), e.Length()}
}
