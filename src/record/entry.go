package record

type OpenRangeStart Time

type Entry struct {
	value   interface{}
	summary Summary
}

func (e Entry) SummaryAsString() string {
	return string(e.summary)
}

func (e Entry) Value() interface{} {
	return e.value
}
