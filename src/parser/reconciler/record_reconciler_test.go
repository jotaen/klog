package reconciler

import (
	. "github.com/jotaen/klog/src"
	"github.com/jotaen/klog/src/parser"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
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
	rs, _ := parser.Parse(original)
	reconciler := NewRecordReconciler(rs, func(r Record) bool {
		return r.Date().ToString() == "2018-01-02"
	})
	require.NotNil(t, reconciler)
	result, err := reconciler.AppendEntry(func(r Record) string { return "2h30m" })
	require.Nil(t, err)
	require.NotNil(t, result)
	require.Equal(t, 150, result.NewRecord.Entries()[2].Duration().InMinutes())
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
`, result.NewText)
}

func TestReconcilerAddsNewlyCreatedEntryAtEndOfFile(t *testing.T) {
	original := `
2018-01-01
    1h`
	rs, _ := parser.Parse(original)
	reconciler := NewRecordReconciler(rs, func(r Record) bool {
		return r.Date().ToString() == "2018-01-01"
	})
	require.NotNil(t, reconciler)
	result, err := reconciler.AppendEntry(func(r Record) string { return "16:00-17:00" })
	require.Nil(t, err)
	assert.Equal(t, `
2018-01-01
    1h
    16:00-17:00
`, result.NewText)
}

func TestReconcilerSkipsIfNoRecordMatches(t *testing.T) {
	original := "2018-01-01\n"
	rs, _ := parser.Parse(original)
	reconciler := NewRecordReconciler(rs, func(r Record) bool { return false })
	require.Nil(t, reconciler)
}

func TestReconcilerRejectsInvalidEntry(t *testing.T) {
	original := "2018-01-01\n"
	rs, _ := parser.Parse(original)
	reconciler := NewRecordReconciler(rs, func(r Record) bool { return true })
	require.NotNil(t, reconciler)
	result, err := reconciler.AppendEntry(func(r Record) string { return "this is not valid entry text" })
	require.Nil(t, result)
	assert.Error(t, err)
}

func TestReconcilerClosesOpenRangeWithNewSummary(t *testing.T) {
	original := `
2018-01-01
    15:00 - ?
`
	rs, _ := parser.Parse(original)
	reconciler := NewRecordReconciler(rs, func(r Record) bool {
		return r.Date().ToString() == "2018-01-01"
	})
	require.NotNil(t, reconciler)
	result, err := reconciler.CloseOpenRange(func(r Record) (Time, EntrySummary) {
		return Ɀ_Time_(15, 22), NewEntrySummary("Finished.")
	})
	require.Nil(t, err)
	assert.Equal(t, `
2018-01-01
    15:00 - 15:22 Finished.
`, result.NewText)
}

func TestReconcilerClosesOpenRangeWithExtendingSummary(t *testing.T) {
	original := `
2018-01-01
    1h
    15:00-??? Will this close? I hope so!?!?
	2m
`
	rs, _ := parser.Parse(original)
	reconciler := NewRecordReconciler(rs, func(r Record) bool {
		return r.Date().ToString() == "2018-01-01"
	})
	require.NotNil(t, reconciler)
	result, err := reconciler.CloseOpenRange(func(r Record) (Time, EntrySummary) {
		return Ɀ_Time_(16, 42), NewEntrySummary("Yes!")
	})
	require.Nil(t, err)
	assert.Equal(t, `
2018-01-01
    1h
    15:00-16:42 Will this close? I hope so!?!? Yes!
	2m
`, result.NewText)
}
