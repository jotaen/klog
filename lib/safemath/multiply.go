package safemath

// Multiply calculates the product of two integers.
// It returns an error if the resulting integer would overflow.
func Multiply(a int, b int) (int, error) {
	err := assertOperandInRange(a, b)
	if err != nil {
		return 0, err
	}

	if b != 0 && abs(a) > MaxInt/abs(b) {
		return 0, ErrOverflow
	}

	return a * b, nil
}
