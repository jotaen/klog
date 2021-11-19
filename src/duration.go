package klog

import (
	"errors"
	"fmt"
	"regexp"
	"strconv"
)

// Duration represents a time span.
type Duration interface {
	InMinutes() int

	// Plus adds up two durations and returns a new duration.
	// It doesn’t alter the original duration object.
	Plus(Duration) Duration

	// Minus subtracts the second from the first duration.
	// It doesn’t alter the original duration object.
	Minus(Duration) Duration

	// ToString serialises the duration. If the duration is negative,
	// the value is preceded by a `-`. E.g. `45m` or `-2h15m`.
	ToString() string

	// ToStringWithSign serialises the duration. In contrast to `ToString`
	// it also precedes positive values with a `+`. If the duration is `0`,
	// no sign will be added. E.g. `-45m` or `0` or `+6h`.
	ToStringWithSign() string
}

func NewDuration(amountHours int, amountMinutes int) Duration {
	return duration(amountHours*60) + duration(amountMinutes)
}

type duration int

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

func (d duration) InMinutes() int {
	return int(d)
}

func (d duration) Plus(additional Duration) Duration {
	return NewDuration(0, d.InMinutes()+additional.InMinutes())
}

func (d duration) Minus(deductible Duration) Duration {
	return NewDuration(0, d.InMinutes()-deductible.InMinutes())
}

func (d duration) ToString() string {
	if d == 0 {
		return "0m"
	}
	hours := abs(int(d) / 60)
	minutes := abs(int(d) % 60)
	result := ""
	if int(d) < 0 {
		result += "-"
	}
	if hours > 0 {
		result += fmt.Sprintf("%dh", hours)
	}
	if minutes > 0 {
		result += fmt.Sprintf("%dm", minutes)
	}
	return result
}

func (d duration) ToStringWithSign() string {
	s := d.ToString()
	if d > 0 {
		return "+" + s
	}
	return s
}

var durationPattern = regexp.MustCompile(`^(-|\+)?((\d+)h)?((\d+)m)?$`)

func NewDurationFromString(hhmm string) (Duration, error) {
	match := durationPattern.FindStringSubmatch(hhmm)
	if match == nil {
		return nil, errors.New("MALFORMED_DURATION")
	}
	sign := 1
	if match[1] == "-" {
		sign = -1
	}
	if match[3] == "" && match[5] == "" {
		return nil, errors.New("MALFORMED_DURATION")
	}
	amountOfHours, _ := strconv.Atoi(match[3])
	amountOfMinutes, _ := strconv.Atoi(match[5])
	if amountOfHours != 0 && amountOfMinutes >= 60 {
		return nil, errors.New("UNREPRESENTABLE_DURATION")
	}
	return NewDuration(sign*amountOfHours, sign*amountOfMinutes), nil
}
