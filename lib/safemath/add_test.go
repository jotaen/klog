package safemath

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestAddNonOverflowing(t *testing.T) {
	for _, x := range []opRes{
		{0, 0, 0},

		{0, 1, 1},
		{0, 42, 42},
		{0, MaxInt, MaxInt},

		{0, -1, -1},
		{0, -42, -42},
		{0, MinInt, MinInt},

		{42, 42, 84},
		{1500000, 1500000, 3000000},

		{MaxInt / 2, MaxInt/2 + 1, MaxInt},
		{MinInt / 2, MinInt/2 - 1, MinInt},

		{MaxInt, -1, MaxInt - 1},
		{MinInt, 1, MinInt + 1},
	} {
		t.Run(fmt.Sprintf("%d %d (original)", x.a, x.b), func(t *testing.T) {
			res, err := Add(x.a, x.b)
			assert.Nil(t, err)
			assert.Equal(t, x.res, res)
		})
		t.Run(fmt.Sprintf("%d %d (commutative)", x.b, x.a), func(t *testing.T) {
			res, err := Add(x.b, x.a)
			assert.Nil(t, err)
			assert.Equal(t, x.res, res)
		})
	}
}

func TestAddOverflowing(t *testing.T) {
	for _, x := range []opErr{
		{1, MaxInt},
		{-1, MinInt},
		{MaxInt, MaxInt},
		{MaxInt/2 + 1, MaxInt/2 + 1},
	} {
		t.Run(fmt.Sprintf("%d %d (original)", x.a, x.b), func(t *testing.T) {
			res, err := Add(x.a, x.b)
			assert.Error(t, err)
			assert.Equal(t, 0, res)
		})
		t.Run(fmt.Sprintf("%d %d (commutative)", x.b, x.a), func(t *testing.T) {
			res, err := Add(x.b, x.a)
			assert.Error(t, err)
			assert.Equal(t, 0, res)
		})
	}
}
