package safemath

func assertOperandInRange(xs ...int) error {
	for _, x := range xs {
		if x < MinInt {
			return ErrOverflow
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
