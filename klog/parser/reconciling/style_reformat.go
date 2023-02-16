package reconciling

import "github.com/jotaen/klog/klog"

// ReformatDirective tells the reconciler whether or in which way the
// date or time value is supposed to be reformatted when serialising it.
type ReformatDirective[T klog.TimeFormat | klog.DateFormat] struct {
	Value T
	mode  int // 0 = Don’t do anything; 1 = Reformat from `Value`; 2 = Apply auto-styling
}

// NoReformat means the time/date value should be taken as is, without touching
// its own existing format.
func NoReformat[T klog.TimeFormat | klog.DateFormat]() ReformatDirective[T] {
	return ReformatDirective[T]{mode: 0}
}

// ReformatExplicitly means that the time/date value should be reformatted
// according to the provided format.
func ReformatExplicitly[T klog.TimeFormat | klog.DateFormat](value T) ReformatDirective[T] {
	return ReformatDirective[T]{Value: value, mode: 1}
}

// ReformatAutoStyle means that the time/date value should be reformatted
// in accordance to what prevalent style the reconciler detects in the file.
// If the style cannot be determined from the file, it falls back to the
// recommended style (as of the file format specification).
func ReformatAutoStyle[T klog.TimeFormat | klog.DateFormat]() ReformatDirective[T] {
	return ReformatDirective[T]{mode: 2}
}

func (r ReformatDirective[T]) apply(autoStyle T, reformat func(T)) {
	if r.mode == 0 {
		return
	}
	format := autoStyle
	if r.mode == 1 {
		format = r.Value
	}
	reformat(format)
}
