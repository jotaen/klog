package safemath

import (
	"errors"
	"math"
)

var (
	OverflowErr = errors.New("overflow")

	// MaxInt represents the largest possible (positive) integer value.
	MaxInt = math.MaxInt
	// MinInt represents the smallest possible (negative) integer value.
	// It doesn’t fully exhaust the theoretical range, to be in line with the
	// MaxInt range, and to allow inverting values without worry.
	MinInt = math.MinInt + 1
)

func assertOperandInRange(xs ...int) error {
	for _, x := range xs {
		if x < MinInt {
			return OverflowErr
		}
	}
	return nil
}

func abs(x int) int {
	if x < 0 {
		return -1 * x
	}
	return x
}
