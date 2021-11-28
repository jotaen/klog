package reconciling

import (
	"errors"
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
	rs, bs, _ := parser.Parse(original)
	reconciler := NewReconciler(rs, bs)
	match := func(r Record) bool {
		return r.Date().ToString() == "2018-01-02"
	}
	require.NotNil(t, reconciler)
	result, err := reconciler.AppendEntry(match, func(r Record) (string, error) { return "2h30m", nil })
	require.Nil(t, err)
	require.NotNil(t, result)
	require.Equal(t, 150, result.Record().Entries()[2].Duration().InMinutes())
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
`, result.FileContents())
}

func TestReconcilerAddsNewlyCreatedEntryAtEndOfFile(t *testing.T) {
	original := `
2018-01-01
    1h`
	rs, bs, _ := parser.Parse(original)

	reconciler := NewReconciler(rs, bs)
	match := func(r Record) bool {
		return r.Date().ToString() == "2018-01-01"
	}
	require.NotNil(t, reconciler)
	result, err := reconciler.AppendEntry(match, func(r Record) (string, error) { return "16:00-17:00", nil })
	require.Nil(t, err)
	assert.Equal(t, `
2018-01-01
    1h
    16:00-17:00
`, result.FileContents())
}

func TestReconcilerSkipsIfNoRecordMatches(t *testing.T) {
	original := "2018-01-01\n"
	rs, bs, _ := parser.Parse(original)
	reconciler := NewReconciler(rs, bs)
	match := func(r Record) bool { return false }
	_, err := reconciler.AppendEntry(match, func(record Record) (string, error) { return "", nil })
	require.ErrorIs(t, err, NotEligibleError{})
}

func TestReconcilerRejectsInvalidEntry(t *testing.T) {
	original := "2018-01-01\n"
	rs, bs, _ := parser.Parse(original)
	reconciler := NewReconciler(rs, bs)
	match := func(r Record) bool { return true }
	require.NotNil(t, reconciler)
	result1, err1 := reconciler.AppendEntry(match, func(r Record) (string, error) { return "this is not valid entry text", nil })
	require.Nil(t, result1)
	assert.Error(t, err1)

	result2, err2 := reconciler.AppendEntry(match, func(r Record) (string, error) { return "", errors.New("Something went wrong") })
	require.Nil(t, result2)
	assert.Error(t, err2)
}

func TestReconcilerClosesOpenRangeWithNewSummary(t *testing.T) {
	original := `
2018-01-01
    15:00 - ?
`
	rs, bs, _ := parser.Parse(original)
	reconciler := NewReconciler(rs, bs)
	match := func(r Record) bool {
		return r.Date().ToString() == "2018-01-01"
	}
	require.NotNil(t, reconciler)
	result, err := reconciler.CloseOpenRange(match, func(r Record) (Time, EntrySummary) {
		return Ɀ_Time_(15, 22), NewEntrySummary("Finished.")
	})
	require.Nil(t, err)
	assert.Equal(t, `
2018-01-01
    15:00 - 15:22 Finished.
`, result.FileContents())
}

func TestReconcilerClosesOpenRangeWithExtendingSummary(t *testing.T) {
	original := `
2018-01-01
    1h
    15:00-??? Will this close? I hope so!?!?
    2m
`
	rs, bs, _ := parser.Parse(original)
	reconciler := NewReconciler(rs, bs)
	match := func(r Record) bool {
		return r.Date().ToString() == "2018-01-01"
	}
	require.NotNil(t, reconciler)
	result, err := reconciler.CloseOpenRange(match, func(r Record) (Time, EntrySummary) {
		return Ɀ_Time_(16, 42), NewEntrySummary("Yes!")
	})
	require.Nil(t, err)
	assert.Equal(t, `
2018-01-01
    1h
    15:00-16:42 Will this close? I hope so!?!? Yes!
    2m
`, result.FileContents())
}

func TestReconcileAddBlockIfOriginalIsEmpty(t *testing.T) {
	rs, bs, _ := parser.Parse("")
	reconciler := NewReconciler(rs, bs)
	date := Ɀ_Date_(3333, 1, 1)
	result, err := reconciler.InsertRecord(date, []InsertableText{
		{"2000-05-05", 0},
	})
	require.Nil(t, err)
	assert.Equal(t, "2000-05-05\n", result.FileContents())
}

func TestReconcileAddBlockToEnd(t *testing.T) {
	original := `
2018-01-01
    1h`
	rs, bs, _ := parser.Parse(original)
	reconciler := NewReconciler(rs, bs)
	date := Ɀ_Date_(2018, 1, 2)
	result, err := reconciler.InsertRecord(date, []InsertableText{
		{"2018-01-02", 0},
	})
	require.Nil(t, err)
	assert.Equal(t, `
2018-01-01
    1h

2018-01-02
`, result.FileContents())
}

func TestReconcileAddBlockToEndWithTrailingNewlines(t *testing.T) {
	original := `
2018-01-01
    1h

`
	rs, bs, _ := parser.Parse(original)
	reconciler := NewReconciler(rs, bs)
	date := Ɀ_Date_(2018, 1, 2)
	result, err := reconciler.InsertRecord(date, []InsertableText{
		{"2018-01-02", 0},
	})
	require.Nil(t, err)
	assert.Equal(t, `
2018-01-01
    1h

2018-01-02

`, result.FileContents())
}

func TestReconcileAddBlockToBeginning(t *testing.T) {
	original := "2018-01-02"
	rs, bs, _ := parser.Parse(original)
	reconciler := NewReconciler(rs, bs)
	date := Ɀ_Date_(2018, 1, 1)
	result, err := reconciler.InsertRecord(date, []InsertableText{
		{"2018-01-01", 0},
	})
	require.Nil(t, err)
	assert.Equal(t, `2018-01-01

2018-01-02`, result.FileContents())
}

func TestReconcileAddBlockToBeginningWithLeadingNewlines(t *testing.T) {
	original := "\n\n2018-01-02"
	rs, bs, _ := parser.Parse(original)
	reconciler := NewReconciler(rs, bs)
	date := Ɀ_Date_(2018, 1, 1)
	result, err := reconciler.InsertRecord(date, []InsertableText{
		{"2018-01-01", 0},
	})
	require.Nil(t, err)
	assert.Equal(t, `2018-01-01



2018-01-02`, result.FileContents())
}

func TestReconcileAddBlockInBetween(t *testing.T) {
	original := `
2018-01-01
    1h

2018-01-03
    3h`
	rs, bs, _ := parser.Parse(original)
	reconciler := NewReconciler(rs, bs)
	date := Ɀ_Date_(2018, 1, 2)
	result, err := reconciler.InsertRecord(date, []InsertableText{
		{"2018-01-02", 0},
		{"This and that", 0},
		{"30m worked", 1},
	})
	require.Nil(t, err)
	assert.Equal(t, `
2018-01-01
    1h

2018-01-02
This and that
    30m worked

2018-01-03
    3h`, result.FileContents())
}
