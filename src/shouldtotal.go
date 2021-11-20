package klog

// ShouldTotal represents the targeted total time of a Record.
type ShouldTotal Duration
type shouldTotal struct {
	Duration
}

func NewShouldTotal(hours int, minutes int) ShouldTotal {
	return shouldTotal{NewDuration(hours, minutes)}
}

func (s shouldTotal) ToString() string {
	return s.Duration.ToString() + "!"
}
