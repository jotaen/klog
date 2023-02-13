package service

import (
	"errors"
	"github.com/jotaen/klog/klog"
	"strconv"
	"strings"
)

// Rounding is an integer divider of 60 that Time values can be rounded to.
type Rounding interface {
	ToInt() int
	ToString() string
}

type rounding int

func (r rounding) ToInt() int {
	return int(r)
}

func (r rounding) ToString() string {
	return strconv.Itoa(r.ToInt()) + "m"
}

// NewRounding creates a Rounding from an integer. For non-allowed
// values, it returns error.
func NewRounding(r int) (Rounding, error) {
	for _, validRounding := range []int{5, 10, 15, 30, 60} {
		if r == validRounding {
			return rounding(r), nil
		}
	}
	return nil, errors.New("INVALID_ROUNDING")
}

// NewRoundingFromString parses a string containing a rounding value.
// The string might be suffixed with `m`. Additionally, it might be `1h`,
// which is equivalent to `60m`.
func NewRoundingFromString(v string) (Rounding, error) {
	r := func() int {
		if v == "1h" {
			return 60
		}
		v = strings.TrimSuffix(v, "m")
		number, err := strconv.Atoi(v)
		if err != nil {
			return -1
		}
		return number
	}()
	return NewRounding(r)
}

// RoundToNearest rounds a time (up or down) to the nearest given rounding multiple.
// E.g., for rounding=5m: 8:03 => 8:05, or for rounding=30m: 15:12 => 15:00
func RoundToNearest(t klog.Time, r Rounding) klog.Time {
	midnightOffset := t.MidnightOffset().InMinutes()
	v := r.ToInt()
	remainder := midnightOffset % v
	uprounder := func() int { // Decide whether to round up the value.
		if remainder >= (v/2 + v%2) {
			return v
		}
		return 0
	}()
	roundedMidnightOffset := midnightOffset - remainder + uprounder

	midnight, _ := klog.NewTime(0, 0)
	roundedTime, err := midnight.Plus(klog.NewDuration(0, roundedMidnightOffset))
	if err != nil {
		// This is the special case where we can’t round up after `23:59>`.
		maxTime, _ := klog.NewTimeTomorrow(23, 59)
		return maxTime
	}
	return roundedTime
}
