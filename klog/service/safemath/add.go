package safemath

// Add calculates the sum of two integers.
// It returns an error if the resulting integer would overflow.
func Add(a int, b int) (int, error) {
	err := assertOperandInRange(a, b)
	if err != nil {
		return 0, err
	}

	if b > 0 {
		if a > MaxInt-b {
			return 0, OverflowErr
		}
	} else {
		if a < MinInt-b {
			return 0, OverflowErr
		}
	}

	return a + b, nil
}
