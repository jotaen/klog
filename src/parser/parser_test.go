package parser

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"klog"
	"testing"
)

func TestParseEmptyDocument(t *testing.T) {
	text := ``
	rs, errs := Parse(text)
	require.Nil(t, errs)
	require.Nil(t, rs)
}

func TestParseBlankDocument(t *testing.T) {
	text := `
 
`
	rs, errs := Parse(text)
	require.Nil(t, errs)
	require.Nil(t, rs)
}

func TestParseMinimalDocument(t *testing.T) {
	text := `2000-01-01`
	rs, errs := Parse(text)
	require.Nil(t, errs)
	require.Len(t, rs, 1)
	assert.Equal(t, src.Ɀ_Date_(2000, 1, 1), rs[0].Date())
}

func TestParseMultipleRecords(t *testing.T) {
	text := `
1999-05-31

1999-06-03 (8h15m!)
Empty
`
	rs, errs := Parse(text)
	require.Nil(t, errs)
	require.Len(t, rs, 2)

	assert.Equal(t, src.Ɀ_Date_(1999, 5, 31), rs[0].Date())
	assert.Equal(t, src.Summary(""), rs[0].Summary())
	assert.Len(t, rs[0].Entries(), 0)

	assert.Equal(t, src.Ɀ_Date_(1999, 6, 3), rs[1].Date())
	assert.Equal(t, src.Summary("Empty"), rs[1].Summary())
	assert.Equal(t, src.NewDuration(8, 15), rs[1].ShouldTotal())
	assert.Len(t, rs[1].Entries(), 0)
}

func TestAcceptTabOrSpacesAsIndentation(t *testing.T) {
	for _, test := range []struct {
		text   string
		expect interface{}
	}{
		{"2000-01-01\n\t8h", nil},
		{"2000-05-31\n  6h", nil},
		{"2000-05-31\n   6h", nil},
		{"2000-05-31\n    6h", nil},
	} {
		rs, errs := Parse(test.text)
		require.Nil(t, errs)
		require.Len(t, rs, 1)
	}
}

func TestParseDocumentSucceedsWithCorrectEntries(t *testing.T) {
	for _, test := range []struct {
		text          string
		expectEntry   interface{}
		expectSummary string
	}{
		{"1234-12-12\n\t5h Some remark", src.NewDuration(5, 0), "Some remark"},
		{"1234-12-12\n\t2h30m", src.NewDuration(2, 30), ""},
		{"1234-12-12\n\t2m", src.NewDuration(0, 2), ""},
		{"1234-12-12\n\t+5h", src.NewDuration(5, 0), ""},
		{"1234-12-12\n\t+2h30m", src.NewDuration(2, 30), ""},
		{"1234-12-12\n\t+2m", src.NewDuration(0, 2), ""},
		{"1234-12-12\n\t-5h", src.NewDuration(-5, -0), ""},
		{"1234-12-12\n\t-2h30m", src.NewDuration(-2, -30), ""},
		{"1234-12-12\n\t-2m", src.NewDuration(-0, -2), ""},
		{"1234-12-12\n\t3:05 - 11:59 Did this and that", src.Ɀ_Range_(src.Ɀ_Time_(3, 5), src.Ɀ_Time_(11, 59)), "Did this and that"},
		{"1234-12-12\n\t<23:30 - 0:10", src.Ɀ_Range_(src.Ɀ_TimeYesterday_(23, 30), src.Ɀ_Time_(0, 10)), ""},
		{"1234-12-12\n\t22:17 - 1:00>", src.Ɀ_Range_(src.Ɀ_Time_(22, 17), src.Ɀ_TimeTomorrow_(1, 00)), ""},
		{"1234-12-12\n\t18:45 - ? Just started something", src.NewOpenRange(src.Ɀ_Time_(18, 45)), "Just started something"},
		{"1234-12-12\n\t<3:12-??????", src.NewOpenRange(src.Ɀ_TimeYesterday_(3, 12)), ""},
	} {
		rs, errs := Parse(test.text)
		require.Nil(t, errs, test.text)
		require.Len(t, rs, 1, test.text)
		require.Len(t, rs[0].Entries(), 1, test.text)
		assert.Equal(t, test.expectEntry, rs[0].Entries()[0].Value(), test.text)
		assert.Equal(t, src.Summary(test.expectSummary), rs[0].Entries()[0].Summary(), test.text)
	}
}

func TestMalformedRecord(t *testing.T) {
	text := `
1999-05-31
	5h30m This and that
Why is there a summary at the end?
`
	rs, errs := Parse(text)
	require.Nil(t, rs)
	require.NotNil(t, errs)
	require.Len(t, errs.Get(), 1)
	assert.Equal(t, Err{ILLEGAL_INDENTATION, 4, 0, 34}, toErr(errs.Get()[0]))
}

func TestReportErrorsInHeadline(t *testing.T) {
	for _, test := range []struct {
		text   string
		expect Err
	}{
		{"Hello 123", Err{INVALID_VALUE, 1, 0, 5}},
		{" 2020-01-01", Err{ILLEGAL_WHITESPACE, 1, 0, 1}},
		{"2020-01-01 (asdf)", Err{UNRECOGNISED_TOKEN, 1, 12, 4}},
		{"2020-01-01 5h30m!", Err{ILLEGAL_SYNTAX, 1, 11, 6}},
		{"2020-01-01 (5h30m!", Err{ILLEGAL_SYNTAX, 1, 18, 1}},
		{"2020-01-01 (", Err{ILLEGAL_SYNTAX, 1, 12, 1}},
		{"2020-01-01 (5h!) foo", Err{ILLEGAL_SYNTAX, 1, 17, 3}},
	} {
		rs, errs := Parse(test.text)
		require.Nil(t, rs)
		require.NotNil(t, errs)
		require.Len(t, errs.Get(), 1)
		assert.Equal(t, test.expect, toErr(errs.Get()[0]), test.text)
	}
}

func TestReportErrorsInSummary(t *testing.T) {
	text := `
2020-01-01
This is a summary that contains
 whitespace at the beginning of the line.
That is not allowed.
`
	rs, errs := Parse(text)
	require.Nil(t, rs)
	require.NotNil(t, errs)
	require.Len(t, errs.Get(), 1)
	assert.Equal(t, Err{INVALID_VALUE, 4, 0, 41}, toErr(errs.Get()[0]))
}

func TestReportErrorsInEntries(t *testing.T) {
	for _, test := range []struct {
		text   string
		expect Err
	}{
		{"2020-01-01\n\t5h1", Err{INVALID_VALUE, 2, 0, 3}},
		{"2020-01-01\n\tasdf Test 123", Err{INVALID_VALUE, 2, 0, 4}},
		{"2020-01-01\n\t15:30", Err{INVALID_VALUE, 2, 5, 1}},
		{"2020-01-01\n\t08:00-", Err{INVALID_VALUE, 2, 6, 1}},
		{"2020-01-01\n\t08:00-asdf", Err{INVALID_VALUE, 2, 6, 4}},
		{"2020-01-01\n\t08:00 - ?asdf", Err{INVALID_VALUE, 2, 9, 4}},
		{"2020-01-01\n\t08:00- ?\n\t09:00 - ?", Err{DUPLICATE_OPEN_RANGE, 3, 0, 9}},
		{"2020-01-01\n\t15:00 - 14:00", Err{ILLEGAL_RANGE, 2, 0, 13}},
		{"2020-01-01\n\t-18:00", Err{INVALID_VALUE, 2, 0, 6}},
		{"2020-01-01\n\t15:30 Foo Bar Baz", Err{INVALID_VALUE, 2, 6, 1}},
	} {
		rs, errs := Parse(test.text)
		require.Nil(t, rs, test.text)
		require.NotNil(t, errs, test.text)
		require.Len(t, errs.Get(), 1, test.text)
		assert.Equal(t, test.expect, toErr(errs.Get()[0]), test.text)
	}
}
