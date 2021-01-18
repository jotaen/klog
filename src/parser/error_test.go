package parser

import . "klog/parser/engine"

type Err struct {
	code string
	line int
	pos  int
	len  int
}

func toErr(e Error) Err {
	return Err{e.Code(), e.Context().LineNumber, e.Position(), e.Length()}
}
