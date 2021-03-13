package parser

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	. "klog"
	"testing"
)

func TestReconcilerAddsNewlyCreatedEntry(t *testing.T) {
	original := `
2018-01-01
    1h

2018-01-02
Hello World
    1h
    1h45m

2018-01-03
    5h
`
	pr, _ := Parse(original)
	reconciled, err := pr.AddEntry(func(rs []Record) (int, string) {
		return 1, "2h30m"
	})
	require.Nil(t, err)
	assert.Equal(t, `
2018-01-01
    1h

2018-01-02
Hello World
    1h
    1h45m
    2h30m

2018-01-03
    5h
`, reconciled)
}

func TestReconcilerAddsNewlyCreatedEntryAtEndOfFile(t *testing.T) {
	original := `
2018-01-01
    1h`
	pr, _ := Parse(original)
	reconciled, err := pr.AddEntry(func(rs []Record) (int, string) {
		return 0, "16:00-17:00"
	})
	require.Nil(t, err)
	assert.Equal(t, `
2018-01-01
    1h
    16:00-17:00
`, reconciled)
}

func TestReconcilerRejectsInvalidIndex(t *testing.T) {
	original := "2018-01-01\n"
	pr, _ := Parse(original)
	reconciled, err := pr.AddEntry(func(rs []Record) (int, string) {
		return 1872, ""
	})
	assert.Equal(t, original, reconciled)
	assert.Error(t, err)
}

func TestReconcilerRejectsNegativeIndex(t *testing.T) {
	original := "2018-01-01\n"
	pr, _ := Parse(original)
	reconciled, err := pr.AddEntry(func(rs []Record) (int, string) {
		return -1, ""
	})
	assert.Equal(t, original, reconciled)
	assert.Error(t, err)
}

func TestReconcilerRejectsInvalidEntry(t *testing.T) {
	original := "2018-01-01\n"
	pr, _ := Parse(original)
	reconciled, err := pr.AddEntry(func(rs []Record) (int, string) {
		return 0, "this is not valid entry text"
	})
	assert.Equal(t, original, reconciled)
	assert.Error(t, err)
}
