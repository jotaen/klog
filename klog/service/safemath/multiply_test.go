package safemath

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestMultiplyNonOverflowing(t *testing.T) {
	for _, x := range []OpRes{
		{0, 0, 0},

		{0, 1, 0},
		{0, 42, 0},
		{0, MaxInt, 0},

		{0, -1, 0},
		{0, -42, 0},
		{0, MinInt, 0},

		{1, MaxInt, MaxInt},
		{-1, MaxInt, MinInt},
		{1, MinInt, MinInt},
		{-1, MinInt, MaxInt},

		{147, 5199, 764253},
		{147, -5199, -764253},
		{-147, -5199, 764253},

		{2, (MaxInt - 1) / 2, MaxInt - 1},
		{2, (MinInt + 1) / 2, MinInt + 1},

		{-1, MinInt, MaxInt},
	} {
		t.Run(fmt.Sprintf("%d %d (original)", x.a, x.b), func(t *testing.T) {
			res, err := Multiply(x.a, x.b)
			require.Nil(t, err)
			assert.Equal(t, x.res, res)
		})
		t.Run(fmt.Sprintf("%d %d (commutative)", x.b, x.a), func(t *testing.T) {
			res, err := Multiply(x.b, x.a)
			require.Nil(t, err)
			assert.Equal(t, x.res, res)
		})
	}
}

func TestMultiplyOverflowing(t *testing.T) {
	for _, x := range []OpErr{
		{1, MaxInt + 1},
		{-1, MaxInt + 1},
		{1, MinInt - 1},
		{-1, MinInt - 1},

		{2, (MaxInt / 2) + 1},
		{2, MaxInt},
		{2, (MinInt / 2) - 1},
		{2, MinInt},

		{MaxInt / 2, MaxInt / 2},

		{MaxInt, MaxInt},
		{MinInt, MinInt},
		{MinInt, MaxInt},
	} {
		t.Run(fmt.Sprintf("%d %d (original)", x.a, x.b), func(t *testing.T) {
			res, err := Multiply(x.a, x.b)
			require.Error(t, err)
			assert.Equal(t, 0, res)
		})
		t.Run(fmt.Sprintf("%d %d (commutative)", x.b, x.a), func(t *testing.T) {
			res, err := Multiply(x.b, x.a)
			require.Error(t, err)
			assert.Equal(t, 0, res)
		})
	}
}
