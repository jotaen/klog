package reconciling

import (
	. "github.com/jotaen/klog/src"
	"github.com/jotaen/klog/src/parser"
	"github.com/jotaen/klog/src/service"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestReconcilerAddsNewEntry(t *testing.T) {
	original := `
2018-01-01
    1h

2018-01-02
Hello World
    1h
    1h45m Multiline...
        ....entry summary

2018-01-03
    5h
`
	rs, _ := parser.Parse(original)
	reconciler := NewReconcilerAtRecord(rs, Ɀ_Date_(2018, 1, 2))
	require.NotNil(t, reconciler)
	result, err := reconciler.AppendEntry("2h30m")
	require.Nil(t, err)
	require.NotNil(t, result)
	require.Equal(t, 150, result.Record.Entries()[2].Duration().InMinutes())
	require.Equal(t, 315, service.Total(result.Record).InMinutes())
	assert.Equal(t, `
2018-01-01
    1h

2018-01-02
Hello World
    1h
    1h45m Multiline...
        ....entry summary
    2h30m

2018-01-03
    5h
`, result.AllSerialised)
}

func TestReconcilerAddsNewlyCreatedEntryAtEndOfFile(t *testing.T) {
	original := "\n2018-01-01\n    1h"
	rs, _ := parser.Parse(original)

	reconciler := NewReconcilerAtRecord(rs, Ɀ_Date_(2018, 1, 1))
	require.NotNil(t, reconciler)
	result, err := reconciler.AppendEntry("16:00-17:00")
	require.Nil(t, err)
	assert.Equal(t, `
2018-01-01
    1h
    16:00-17:00
`, result.AllSerialised)
}

func TestReconcilerRespectsIndentationStyle(t *testing.T) {
	for _, x := range []struct {
		original string
		expected string
	}{
		{"1444-10-09\n\t1h", "1444-10-09\n\t1h\n\t30m\n"},
		{"1444-10-09\n  1h", "1444-10-09\n  1h\n  30m\n"},
		{"1444-10-09\n   1h", "1444-10-09\n   1h\n   30m\n"},
		{"1444-10-09\n    1h", "1444-10-09\n    1h\n    30m\n"},
		{"1444-10-08\n  3h\n\n1444-10-09\n\t1h", "1444-10-08\n  3h\n\n1444-10-09\n\t1h\n\t30m\n"},
	} {
		rs, _ := parser.Parse(x.original)
		reconciler := NewReconcilerAtRecord(rs, Ɀ_Date_(1444, 10, 9))
		result, err := reconciler.AppendEntry("30m")
		require.Nil(t, err)
		assert.Equal(t, x.expected, result.AllSerialised)
	}
}

func TestReconcilerRespectsLineEndingStyle(t *testing.T) {
	for _, x := range []struct {
		original string
		expected string
	}{
		{"1444-10-09\r\n\t1h", "1444-10-09\r\n\t1h\r\n\t30m\r\n"},
		{"1444-10-09\r\n\t1h\n\t2h", "1444-10-09\r\n\t1h\n\t2h\r\n\t30m\r\n"},
	} {
		rs, _ := parser.Parse(x.original)
		reconciler := NewReconcilerAtRecord(rs, Ɀ_Date_(1444, 10, 9))
		result, err := reconciler.AppendEntry("30m")
		require.Nil(t, err)
		assert.Equal(t, x.expected, result.AllSerialised)
	}
}

func TestReconcilerSkipsIfNoRecordMatches(t *testing.T) {
	original := "2018-01-01\n"
	rs, _ := parser.Parse(original)
	reconciler := NewReconcilerAtRecord(rs, Ɀ_Date_(9999, 12, 31))
	require.Nil(t, reconciler)
}

func TestReconcilerRejectsInvalidEntry(t *testing.T) {
	original := "2018-01-01\n"
	rs, _ := parser.Parse(original)
	reconciler := NewReconcilerAtRecord(rs, Ɀ_Date_(2018, 1, 1))
	require.NotNil(t, reconciler)
	result, err := reconciler.AppendEntry("this is not valid entry text")
	require.Nil(t, result)
	assert.Error(t, err)
}

func TestReconcilerStartsOpenRange(t *testing.T) {
	original := `
2018-01-01
	5h22m
`
	rs, _ := parser.Parse(original)
	reconciler := NewReconcilerAtRecord(rs, Ɀ_Date_(2018, 1, 1))
	require.NotNil(t, reconciler)
	result, err := reconciler.StartOpenRange(Ɀ_Time_(8, 3), "")
	require.Nil(t, err)
	assert.Equal(t, `
2018-01-01
	5h22m
	8:03 - ?
`, result.AllSerialised)
}

func TestReconcilerStartsOpenRangeWithNewSummary(t *testing.T) {
	original := `
2018-01-01
	5h22m
`
	rs, _ := parser.Parse(original)
	reconciler := NewReconcilerAtRecord(rs, Ɀ_Date_(2018, 1, 1))
	require.NotNil(t, reconciler)
	result, err := reconciler.StartOpenRange(Ɀ_Time_(8, 3), "Started!")
	require.Nil(t, err)
	assert.Equal(t, `
2018-01-01
	5h22m
	8:03 - ? Started!
`, result.AllSerialised)
}

func TestReconcilerStartsOpenRangeWithStyle(t *testing.T) {
	original := `
2018-01-01
	2:00am-3:00am
`
	rs, _ := parser.Parse(original)
	reconciler := NewReconcilerAtRecord(rs, Ɀ_Date_(2018, 1, 1))
	require.NotNil(t, reconciler)
	result, err := reconciler.StartOpenRange(Ɀ_Time_(8, 3), "")
	require.Nil(t, err)
	// Conforms to both am/pm and spaces around dash
	assert.Equal(t, `
2018-01-01
	2:00am-3:00am
	8:03am-?
`, result.AllSerialised)
}

func TestReconcilerStartsOpenRangeWithStyleFromOtherRecord(t *testing.T) {
	original := `
2018-01-01
  2:00am-3:00am

2018-01-02
`
	rs, _ := parser.Parse(original)
	reconciler := NewReconcilerAtRecord(rs, Ɀ_Date_(2018, 1, 2))
	require.NotNil(t, reconciler)
	result, err := reconciler.StartOpenRange(Ɀ_Time_(8, 3), "")
	require.Nil(t, err)
	// Conforms to both am/pm and spaces around dash
	assert.Equal(t, `
2018-01-01
  2:00am-3:00am

2018-01-02
  8:03am-?
`, result.AllSerialised)
}

func TestReconcilerClosesOpenRangeWithNewSummary(t *testing.T) {
	original := `
2018-01-01
    15:00 - ?
`
	rs, _ := parser.Parse(original)
	reconciler := NewReconcilerAtRecord(rs, Ɀ_Date_(2018, 1, 1))
	require.NotNil(t, reconciler)
	result, err := reconciler.CloseOpenRange(Ɀ_Time_(15, 22), "Finished.")
	require.Nil(t, err)
	assert.Equal(t, `
2018-01-01
    15:00 - 15:22 Finished.
`, result.AllSerialised)
}

func TestReconcilerClosesOpenRangeWithStyle(t *testing.T) {
	original := `
2010-04-27
    3:00pm - ??
`
	rs, _ := parser.Parse(original)
	reconciler := NewReconcilerAtRecord(rs, Ɀ_Date_(2010, 4, 27))
	require.NotNil(t, reconciler)
	result, err := reconciler.CloseOpenRange(Ɀ_Time_(15, 30), "")
	require.Nil(t, err)
	assert.Equal(t, `
2010-04-27
    3:00pm - 3:30pm
`, result.AllSerialised)
}

func TestReconcilerClosesOpenRangeWithExtendingSummary(t *testing.T) {
	original := `
2018-01-01
    1h Multiline...
        ...entry summary
    15:00-??? Will this close?
        I hope so.
    2m
`
	rs, _ := parser.Parse(original)
	reconciler := NewReconcilerAtRecord(rs, Ɀ_Date_(2018, 1, 1))
	require.NotNil(t, reconciler)
	result, err := reconciler.CloseOpenRange(Ɀ_Time_(16, 42), "Yes!")
	require.Nil(t, err)
	assert.Equal(t, `
2018-01-01
    1h Multiline...
        ...entry summary
    15:00-16:42 Will this close?
        I hope so. Yes!
    2m
`, result.AllSerialised)
}

func TestReconcileAddRecordIfOriginalIsEmpty(t *testing.T) {
	rs, _ := parser.Parse("")
	reconciler := NewReconcilerAtNewRecord(rs, Ɀ_Date_(2000, 5, 5), nil)
	result, err := reconciler.MakeResult()
	require.Nil(t, err)
	assert.Equal(t, "2000-05-05\n", result.AllSerialised)
	assert.Equal(t, "2000-05-05", result.Record.Date().ToString())
}

func TestReconcileAddRecordIfOriginalContainsOneRecord(t *testing.T) {
	rs, _ := parser.Parse("1999-12-31")
	reconciler := NewReconcilerAtNewRecord(rs, Ɀ_Date_(2000, 2, 1), nil)
	result, err := reconciler.MakeResult()
	require.Nil(t, err)
	assert.Equal(t, "1999-12-31\n\n2000-02-01\n", result.AllSerialised)
	assert.Equal(t, "2000-02-01", result.Record.Date().ToString())
}

func TestReconcileNewRecordFromEmptyFile(t *testing.T) {
	for _, x := range []struct {
		original string
	}{
		{""},
		{"\n"},
		{"\n\n"},
		{"\n\n\t\n"},
		{"\n\n     \t\n \t     \n  "},
	} {
		rs, _ := parser.Parse(x.original)
		reconciler := NewReconcilerAtNewRecord(rs, Ɀ_Date_(1995, 3, 17), nil)
		result, err := reconciler.MakeResult()
		require.Nil(t, err)
		assert.Equal(t, "1995-03-17\n", result.AllSerialised)
	}
}

func TestReconcilePrependNewRecord(t *testing.T) {
	for _, x := range []struct {
		original string
		expected string
	}{
		{"2018-01-02", "2018-01-01\n\n2018-01-02"},
		{"2018-01-02", "2018-01-01\n\n2018-01-02"},
		{"2018-01-02\n\t1h", "2018-01-01\n\n2018-01-02\n\t1h"},
		{"2018-01-02\n\n", "2018-01-01\n\n2018-01-02\n\n"},
		{"\n\n2018-01-02\n", "2018-01-01\n\n\n\n2018-01-02\n"},
	} {
		rs, _ := parser.Parse(x.original)
		reconciler := NewReconcilerAtNewRecord(rs, Ɀ_Date_(2018, 1, 1), nil)
		result, err := reconciler.MakeResult()
		require.Nil(t, err)
		assert.Equal(t, x.expected, result.AllSerialised)
	}
}

func TestReconcileAppendNewRecord(t *testing.T) {
	for _, x := range []struct {
		original string
		expected string
	}{
		{"2018-01-01", "2018-01-01\n\n2019-01-01\n"},
		{"2018-01-01\n\n", "2018-01-01\n\n2019-01-01\n\n"},
		{"\n\n2018-01-01\n", "\n\n2018-01-01\n\n2019-01-01\n"},
	} {
		rs, _ := parser.Parse(x.original)
		reconciler := NewReconcilerAtNewRecord(rs, Ɀ_Date_(2019, 1, 1), nil)
		result, err := reconciler.MakeResult()
		require.Nil(t, err)
		assert.Equal(t, x.expected, result.AllSerialised)
	}
}

func TestReconcileAddBlockInBetween(t *testing.T) {
	for _, x := range []struct {
		original string
		expected string
	}{
		{"2018-01-01\n\n2018-01-03", "2018-01-01\n\n2018-01-02\n\n2018-01-03"},
		{"2018-01-01\n\n\n2018-01-03", "2018-01-01\n\n2018-01-02\n\n\n2018-01-03"},
		{"2018-01-02\n\t1h\n\n2018-01-03", "2018-01-02\n\t1h\n\n2018-01-02\n\n2018-01-03"},
	} {
		rs, _ := parser.Parse(x.original)
		reconciler := NewReconcilerAtNewRecord(rs, Ɀ_Date_(2018, 1, 2), nil)
		result, err := reconciler.MakeResult()
		require.Nil(t, err)
		assert.Equal(t, x.expected, result.AllSerialised)
	}
}

func TestReconcileAddRecordWithShouldTotal(t *testing.T) {
	original := `
2018-01-01
    1h`
	rs, _ := parser.Parse(original)
	reconciler := NewReconcilerAtNewRecord(rs, Ɀ_Date_(2018, 1, 2), NewShouldTotal(5, 31))
	result, err := reconciler.MakeResult()
	require.Nil(t, err)
	assert.Equal(t, `
2018-01-01
    1h

2018-01-02 (5h31m!)
`, result.AllSerialised)
	assert.Equal(t, NewShouldTotal(5, 31), result.Record.ShouldTotal())
}

func TestReconcileRespectsExistingStylePref(t *testing.T) {
	for _, x := range []struct {
		original string
		expected string
	}{
		{"3145/06/15\n", "3145/06/15\n\n3145/06/16\n"},
		{"3145/06/14\n\n3145/06/15\n\n3145-06-15\n", "3145/06/14\n\n3145/06/15\n\n3145-06-15\n\n3145/06/16\n"},
	} {
		rs, _ := parser.Parse(x.original)
		reconciler := NewReconcilerAtNewRecord(rs, Ɀ_Date_(3145, 6, 16), nil)
		result, err := reconciler.MakeResult()
		require.Nil(t, err)
		assert.Equal(t, x.expected, result.AllSerialised)
	}
}
