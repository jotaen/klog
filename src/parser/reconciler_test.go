package parser

import (
	"errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	. "klog"
	"klog/parser/parsing"
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
	reconciler, err := NewRecordReconciler(
		pr,
		nil,
		func(r Record) bool { return r.Date().ToString() == "2018-01-02" },
	)
	require.Nil(t, err)
	newRecord, reconciled, err := reconciler.AppendEntry(func(r Record) string { return "2h30m" })
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
	reconciler, err := NewRecordReconciler(
		pr,
		nil,
		func(r Record) bool { return r.Date().ToString() == "2018-01-01" },
	)
	require.Nil(t, err)
	_, reconciled, err := reconciler.AppendEntry(func(r Record) string { return "16:00-17:00" })
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
	reconciler, err := NewRecordReconciler(
		pr,
		errors.New("No such record"),
		func(r Record) bool { return false },
	)
	require.Nil(t, reconciler)
	assert.EqualError(t, err, "No such record")
}

func TestReconcilerRejectsInvalidEntry(t *testing.T) {
	original := "2018-01-01\n"
	pr, _ := Parse(original)
	reconciler, err := NewRecordReconciler(
		pr,
		errors.New("No such record"),
		func(r Record) bool { return true },
	)
	require.Nil(t, err)
	newRecord, reconciled, err := reconciler.AppendEntry(func(r Record) string { return "this is not valid entry text" })
	assert.Equal(t, "", reconciled)
	assert.Nil(t, newRecord)
	assert.Error(t, err)
}

func TestReconcilerClosesOpenRange(t *testing.T) {
	original := `
2018-01-01
    1h
    15:00-??? Will this close? I hope so!?!?
	2m
`
	pr, _ := Parse(original)
	reconciler, err := NewRecordReconciler(
		pr,
		nil,
		func(r Record) bool { return r.Date().ToString() == "2018-01-01" },
	)
	require.Nil(t, err)
	_, reconciled, err := reconciler.CloseOpenRange(func(r Record) Time { return Ɀ_Time_(16, 42) })
	require.Nil(t, err)
	assert.Equal(t, `
2018-01-01
    1h
    15:00-16:42 Will this close? I hope so!?!?
	2m
`, reconciled)
}

func TestReconcileAddBlockIfOriginalIsEmpty(t *testing.T) {
	pr, _ := Parse("")
	reconciler, _ := NewBlockReconciler(pr, func(Record, Record) bool {
		return false
	})
	_, reconciled, err := reconciler.AddBlock([]parsing.Text{
		{"2000-05-05", 0},
	})
	require.Nil(t, err)
	assert.Equal(t, "2000-05-05\n", reconciled)
}

func TestReconcileAddBlockToEnd(t *testing.T) {
	original := `
2018-01-01
    1h`
	pr, _ := Parse(original)
	reconciler, _ := NewBlockReconciler(
		pr,
		func(Record, Record) bool { return false },
	)
	_, reconciled, err := reconciler.AddBlock([]parsing.Text{
		{"2018-01-02", 0},
	})
	require.Nil(t, err)
	assert.Equal(t, `
2018-01-01
    1h

2018-01-02
`, reconciled)
}

func TestReconcileAddBlockInBetween(t *testing.T) {
	original := `
2018-01-01
    1h

2018-01-03
    3h`
	pr, _ := Parse(original)
	date := Ɀ_Date_(2018, 1, 2)
	reconciler, _ := NewBlockReconciler(pr, func(r1 Record, r2 Record) bool {
		return date.IsAfterOrEqual(r1.Date()) && r2.Date().IsAfterOrEqual(date)
	})
	_, reconciled, err := reconciler.AddBlock([]parsing.Text{
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
    3h`, reconciled)
}
