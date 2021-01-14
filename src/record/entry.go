package record

import "klog/datetime"

type Entry interface {
	Total() datetime.Duration
	Summary()  string
	IsDuration() bool
	IsTimeRange() bool
	IsOpenRange() bool

	val() interface{}
}

type entry struct {
	value   interface{}
	summary string
}

func (e entry) Total() datetime.Duration {
	switch x := e.value.(type) {
	case datetime.Duration:
		return x
	case datetime.TimeRange:
		return x.Duration()
	}
	return nil
}

func (e entry) Summary() string {
	return e.summary
}

func (e entry) IsDuration() bool {
	_, ok := e.value.(datetime.Duration)
	return ok
}

func (e entry) IsTimeRange() bool {
	_, ok := e.value.(datetime.TimeRange)
	return ok
}

func (e entry) IsOpenRange() bool {
	_, ok := e.value.(datetime.Time)
	return ok
}

func (e entry) val() interface{} {
	return e.value
}
