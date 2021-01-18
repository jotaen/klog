package record

type Entry struct {
	value   interface{}
	summary Summary
}

func (e *Entry) Summary() Summary {
	return e.summary
}

func (e *Entry) Value() interface{} {
	return e.value
}

func (e *Entry) ToString() string {
	switch x := e.Value().(type) {
	case Range:
		return x.ToString()
	case Duration:
		return x.ToString()
	case OpenRange:
		return x.ToString()
	}
	panic("Incomplete switch statement")
}

func (e *Entry) Duration() Duration {
	switch x := e.Value().(type) {
	case Range:
		return x.Duration()
	case Duration:
		return x
	case OpenRange:
		return NewDuration(0, 0)
	}
	panic("Incomplete switch statement")
}
