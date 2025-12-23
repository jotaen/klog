package safemath

import (
	"errors"
	"math"
)

var (
	ErrOverflow = errors.New("overflow")

	// MaxInt represents the largest possible (positive) integer value.
	MaxInt = math.MaxInt
	// MinInt represents the smallest possible (negative) integer value.
	// It doesn’t fully exhaust the theoretical range, to be in line with the
	// MaxInt range, and to allow inverting values without worry.
	MinInt = math.MinInt + 1
)
