package record

type Entry interface {
	Value() interface{}
	Summary() string
}

type OpenRangeStart Time

type entry struct {
	value   interface{}
	summary string
}

func (e entry) Summary() string {
	return e.summary
}

func (e entry) Value() interface{} {
	return e.value
}
