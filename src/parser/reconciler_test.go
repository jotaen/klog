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
	newRecord, reconciled, err := pr.AddEntry(
		"",
		func(r Record) bool { return r.Date().ToString() == "2018-01-02" },
		func(r Record) string { return "2h30m" },
	)
	require.Nil(t, err)
	require.Equal(t, 150, newRecord.Entries()[2].Duration().InMinutes())
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
	_, reconciled, err := pr.AddEntry(
		"",
		func(r Record) bool { return r.Date().ToString() == "2018-01-01" },
		func(r Record) string { return "16:00-17:00" },
	)
	require.Nil(t, err)
	assert.Equal(t, `
2018-01-01
    1h
    16:00-17:00
`, reconciled)
}

func TestReconcilerSkipsIfNoRecordMatches(t *testing.T) {
	original := "2018-01-01\n"
	pr, _ := Parse(original)
	newRecord, reconciled, err := pr.AddEntry(
		"No such record",
		func(r Record) bool { return false },
		func(r Record) string { return "1h" },
	)
	assert.Equal(t, original, reconciled)
	assert.Nil(t, newRecord)
	assert.EqualError(t, err, "No such record")
}

func TestReconcilerRejectsInvalidEntry(t *testing.T) {
	original := "2018-01-01\n"
	pr, _ := Parse(original)
	newRecord, reconciled, err := pr.AddEntry(
		"",
		func(r Record) bool { return true },
		func(r Record) string { return "this is not valid entry text" },
	)
	assert.Equal(t, original, reconciled)
	assert.Nil(t, newRecord)
	assert.Error(t, err)
}
