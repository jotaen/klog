package record

type OpenRangeStart Time

type Entry struct {
	value   interface{}
	summary Summary
}

func (e Entry) Summary() Summary {
	return e.summary
}

func (e Entry) Value() interface{} {
	return e.value
}
