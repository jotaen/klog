package parser

import (
	. "klog/parser/parsing"
	"reflect"
	"runtime"
	"strings"
)

type Err struct {
	id   string
	line int
	pos  int
	len  int
}

func id(fn interface{}) string {
	fullyQualifiedName := runtime.FuncForPC(reflect.ValueOf(fn).Pointer()).Name()
	return fullyQualifiedName[strings.Index(fullyQualifiedName, ".")+1:]
}

func toErr(e Error) Err {
	return Err{e.Code(), e.Context().LineNumber, e.Position(), e.Length()}
}
