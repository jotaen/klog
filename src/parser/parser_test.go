package parser

import (
	. "github.com/jotaen/klog/src"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestParseMinimalDocument(t *testing.T) {
	text := `2000-01-01`
	rs, errs := Parse(text)
	require.Nil(t, errs)
	require.Len(t, rs, 1)
	assert.Equal(t, Ɀ_Date_(2000, 1, 1), rs[0].Date())
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

	assert.Equal(t, Ɀ_Date_(1999, 5, 31), rs[0].Date())
	assert.Equal(t, Ɀ_RecordSummary_(), rs[0].Summary())
	assert.Len(t, rs[0].Entries(), 0)

	assert.Equal(t, Ɀ_Date_(1999, 6, 3), rs[1].Date())
	assert.Equal(t, Ɀ_RecordSummary_("Empty"), rs[1].Summary())
	assert.Equal(t, NewDuration(8, 15).InMinutes(), rs[1].ShouldTotal().InMinutes())
	assert.Len(t, rs[1].Entries(), 0)
}

func TestParseEmptyOrBlankDocument(t *testing.T) {
	for _, text := range []string{
		"",
		"    ",
		"\n\n\n\n\n",
		"\n\t     \n \n         ",
	} {
		rs, errs := Parse(text)
		require.Nil(t, errs)
		require.Len(t, rs, 0)
	}
}

func TestParseWindowsAndUnixLineEndings(t *testing.T) {
	text := "2000-01-01\r\n\r\n2000-01-02\n\n2000-01-03"
	rs, errs := Parse(text)
	require.Nil(t, errs)
	require.Len(t, rs, 3)
}

func TestParseMultipleRecordsWhenBlankLineContainsWhitespace(t *testing.T) {
	text := "2018-01-01\n    1h\n" + "    \n" + "2019-01-01\n"
	rs, errs := Parse(text)
	require.Nil(t, errs)
	require.Len(t, rs, 2)
}

func TestParseAlternativeFormatting(t *testing.T) {
	text := `
1999/05/31
	8:00-13:00

1999-05-31
	8:00am-1:00pm
`
	rs, errs := Parse(text)
	require.Nil(t, errs)
	require.Len(t, rs, 2)

	assert.True(t, rs[0].Date().IsEqualTo(rs[1].Date()))
	assert.Equal(t, rs[0].Entries()[0].Duration(), rs[1].Entries()[0].Duration())
}

func TestAcceptTabOrSpacesAsIndentation(t *testing.T) {
	for _, test := range []struct {
		text   string
		expect interface{}
	}{
		{"2000-01-01\n\t8h", nil},
		{"2000-01-01\n\t8h\n\t15m", nil},
		{"2000-05-31\n  6h", nil},
		{"2000-05-31\n  6h\n  20m", nil},
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
		expectSummary EntrySummary
	}{
		{"1234-12-12\n\t5h Some remark", NewDuration(5, 0), Ɀ_EntrySummary_("Some remark")},
		{"1234-12-12\n\t5h    Some remark", NewDuration(5, 0), Ɀ_EntrySummary_("   Some remark")},
		{"1234-12-12\n\t5h\tSome remark", NewDuration(5, 0), Ɀ_EntrySummary_("Some remark")},
		{"1234-12-12\n\t2h30m", NewDuration(2, 30), nil},
		{"1234-12-12\n\t2h30m ", NewDuration(2, 30), nil},
		{"1234-12-12\n\t2m", NewDuration(0, 2), nil},
		{"1234-12-12\n\t+5h", NewDuration(5, 0), nil},
		{"1234-12-12\n\t+2h30m", NewDuration(2, 30), nil},
		{"1234-12-12\n\t+2m", NewDuration(0, 2), nil},
		{"1234-12-12\n\t-5h", NewDuration(-5, -0), nil},
		{"1234-12-12\n\t-2h30m", NewDuration(-2, -30), nil},
		{"1234-12-12\n\t-2m", NewDuration(-0, -2), nil},
		{"1234-12-12\n\t3:05 - 11:59", Ɀ_Range_(Ɀ_Time_(3, 5), Ɀ_Time_(11, 59)), nil},
		{"1234-12-12\n\t3:05 - 11:59 ", Ɀ_Range_(Ɀ_Time_(3, 5), Ɀ_Time_(11, 59)), nil},
		{"1234-12-12\n\t3:05 - 11:59 Did this and that", Ɀ_Range_(Ɀ_Time_(3, 5), Ɀ_Time_(11, 59)), Ɀ_EntrySummary_("Did this and that")},
		{"1234-12-12\n\t3:05 - 11:59   Foo", Ɀ_Range_(Ɀ_Time_(3, 5), Ɀ_Time_(11, 59)), Ɀ_EntrySummary_("  Foo")},
		{"1234-12-12\n\t3:05 - 11:59\tFoo", Ɀ_Range_(Ɀ_Time_(3, 5), Ɀ_Time_(11, 59)), Ɀ_EntrySummary_("Foo")},
		{"1234-12-12\n\t22:00 - 24:00\tFoo", Ɀ_Range_(Ɀ_Time_(22, 0), Ɀ_TimeTomorrow_(0, 0)), Ɀ_EntrySummary_("Foo")},
		{"1234-12-12\n\t<22:00 - <24:00\tFoo", Ɀ_Range_(Ɀ_TimeYesterday_(22, 0), Ɀ_Time_(0, 0)), Ɀ_EntrySummary_("Foo")},
		{"1234-12-12\n\t9:00am - 1:43pm", Ɀ_Range_(Ɀ_IsAmPm_(Ɀ_Time_(9, 00)), Ɀ_IsAmPm_(Ɀ_Time_(13, 43))), nil},
		{"1234-12-12\n\t9:00am-1:43pm", Ɀ_Range_(Ɀ_IsAmPm_(Ɀ_Time_(9, 00)), Ɀ_IsAmPm_(Ɀ_Time_(13, 43))), nil},
		{"1234-12-12\n\t9:00am-8:12am> Things", Ɀ_Range_(Ɀ_IsAmPm_(Ɀ_Time_(9, 00)), Ɀ_IsAmPm_(Ɀ_TimeTomorrow_(8, 12))), Ɀ_EntrySummary_("Things")},
		{"1234-12-12\n\t9:00am-9:05 Mixed styles", Ɀ_Range_(Ɀ_IsAmPm_(Ɀ_Time_(9, 00)), Ɀ_Time_(9, 05)), Ɀ_EntrySummary_("Mixed styles")},
		{"1234-12-12\n\t<23:30 - 0:10", Ɀ_Range_(Ɀ_TimeYesterday_(23, 30), Ɀ_Time_(0, 10)), nil},
		{"1234-12-12\n\t22:17 - 1:00>", Ɀ_Range_(Ɀ_Time_(22, 17), Ɀ_TimeTomorrow_(1, 00)), nil},
		{"1234-12-12\n\t<23:00-1:00>", Ɀ_Range_(Ɀ_TimeYesterday_(23, 00), Ɀ_TimeTomorrow_(1, 00)), nil},
		{"1234-12-12\n\t<23:00-<23:10", Ɀ_Range_(Ɀ_TimeYesterday_(23, 00), Ɀ_TimeYesterday_(23, 10)), nil},
		{"1234-12-12\n\t12:01>-13:59>", Ɀ_Range_(Ɀ_TimeTomorrow_(12, 01), Ɀ_TimeTomorrow_(13, 59)), nil},
		{"1234-12-12\n\t12:01 - ?", NewOpenRange(Ɀ_Time_(12, 1)), nil},
		{"1234-12-12\n\t12:01 - ? ", NewOpenRange(Ɀ_Time_(12, 1)), nil},
		{"1234-12-12\n\t18:45 - ? Just started something", NewOpenRange(Ɀ_Time_(18, 45)), Ɀ_EntrySummary_("Just started something")},
		{"1234-12-12\n\t18:45 - ?\tFoo", NewOpenRange(Ɀ_Time_(18, 45)), Ɀ_EntrySummary_("Foo")},
		{"1234-12-12\n\t18:45 - ???       ASDF", NewOpenRange(Ɀ_Time_(18, 45)), Ɀ_EntrySummary_("      ASDF")},
		{"1234-12-12\n\t<3:12-??????", NewOpenRange(Ɀ_TimeYesterday_(3, 12)), nil},
	} {
		rs, errs := Parse(test.text)
		require.Nil(t, errs, test.text)
		require.Len(t, rs, 1, test.text)
		require.Len(t, rs[0].Entries(), 1, test.text)
		value := rs[0].Entries()[0].Unbox(
			func(r Range) interface{} { return r },
			func(d Duration) interface{} { return d },
			func(o OpenRange) interface{} { return o },
		)
		assert.Equal(t, test.expectEntry, value, test.text)
		assert.Equal(t, test.expectSummary, rs[0].Entries()[0].Summary(), test.text)
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
	require.Len(t, errs, 1)
	assert.Equal(t, ErrorIllegalIndentation().toErrData(4, 0, 34), toErrData(errs[0]))
}

func TestReportErrorsInHeadline(t *testing.T) {
	for _, test := range []struct {
		text   string
		expect errData
	}{
		{"Hello 123", ErrorInvalidDate().toErrData(1, 0, 5)},
		{" 2020-01-01", ErrorIllegalIndentation().toErrData(1, 0, 11)},
		{"   2020-01-01", ErrorIllegalIndentation().toErrData(1, 0, 13)},
		{"2020-01-01 ()", ErrorMalformedPropertiesSyntax().toErrData(1, 12, 1)},
		{"2020-01-01 (asdf)", ErrorUnrecognisedProperty().toErrData(1, 12, 4)},
		{"2020-01-01 (asdf!)", ErrorMalformedShouldTotal().toErrData(1, 12, 4)},
		{"2020-01-01 5h30m!", ErrorUnrecognisedTextInHeadline().toErrData(1, 11, 6)},
		{"2020-01-01 (5h30m!", ErrorMalformedPropertiesSyntax().toErrData(1, 18, 1)},
		{"2020-01-01 (", ErrorMalformedPropertiesSyntax().toErrData(1, 12, 1)},
		{"2020-01-01 (5h!) foo", ErrorUnrecognisedTextInHeadline().toErrData(1, 17, 3)},
		{"2020-01-01 (5h! asdf)", ErrorUnrecognisedProperty().toErrData(1, 16, 4)},
		{"2020-01-01 (5h!!!)", ErrorUnrecognisedProperty().toErrData(1, 15, 2)},
	} {
		rs, errs := Parse(test.text)
		require.Nil(t, rs)
		require.NotNil(t, errs)
		require.Len(t, errs, 1)
		assert.Equal(t, test.expect, toErrData(errs[0]), test.text)
	}
}

func TestReportErrorsInSummary(t *testing.T) {
	text := `
2020-01-01
This is a summary that contains
 whitespace at the beginning of the line.
That is not allowed.
 Other kinds of blank characters are not allowed there neither.
 And neither are fake blank lines:
    
End.
`
	rs, errs := Parse(text)
	require.Nil(t, rs)
	require.NotNil(t, errs)
	require.Len(t, errs, 4)
	assert.Equal(t, ErrorMalformedSummary().toErrData(4, 0, 41), toErrData(errs[0]))
	assert.Equal(t, ErrorMalformedSummary().toErrData(6, 0, 63), toErrData(errs[1]))
	assert.Equal(t, ErrorMalformedSummary().toErrData(7, 0, 34), toErrData(errs[2]))
	assert.Equal(t, ErrorMalformedSummary().toErrData(8, 0, 4), toErrData(errs[3]))
}

func TestReportErrorsIfIndentationIsIncorrect(t *testing.T) {
	for _, test := range []struct {
		text   string
		expect errData
	}{
		// To few characters (that’s actually a malformed summary, though):
		{"2020-01-01\n 8h", ErrorMalformedSummary().toErrData(2, 0, 3)},

		// Not exactly one indentation level:
		{"2020-01-01\n\t 8h", ErrorIllegalIndentation().toErrData(2, 0, 4)},
		{"2020-01-01\n\t\t8h", ErrorIllegalIndentation().toErrData(2, 0, 4)},
		{"2020-01-01\n     8h", ErrorIllegalIndentation().toErrData(2, 0, 7)},

		// Mixed styles within one record:
		{"2020-01-01\n    8h\n\t2h", ErrorIllegalIndentation().toErrData(3, 0, 3)},
		{"2020-01-01\n  8h\n\t2h", ErrorIllegalIndentation().toErrData(3, 0, 3)},
		{"2020-01-01\n\t8h\n    2h", ErrorIllegalIndentation().toErrData(3, 0, 6)},
		{"2020-01-01\n\t8h\n  2h", ErrorIllegalIndentation().toErrData(3, 0, 4)},
		{"2020-01-01\n  8h\n    2h", ErrorIllegalIndentation().toErrData(3, 0, 6)},
		{"2020-01-01\n    8h\n  2h", ErrorIllegalIndentation().toErrData(3, 0, 4)},
		{"2020-01-01\n  8h\n   2h", ErrorIllegalIndentation().toErrData(3, 0, 5)},
		{"2020-01-01\n  8h\n  2h\n\t1h2m", ErrorIllegalIndentation().toErrData(4, 0, 5)},
	} {
		rs, errs := Parse(test.text)
		require.Nil(t, rs, test.text)
		require.NotNil(t, errs, test.text)
		require.Len(t, errs, 1, test.text)
		assert.Equal(t, test.expect, toErrData(errs[0]), test.text)
	}
}

func TestAcceptMixingIndentationStylesAcrossRecords(t *testing.T) {
	text := `
2020-01-01
  4h This is two spaces
  2h So is this

2020-01-02
    6h This is 4 spaces

2020-01-03
	12m This is a tab
`
	rs, errs := Parse(text)
	require.Nil(t, errs)
	require.Len(t, rs, 3)
}

func TestReportErrorsInEntries(t *testing.T) {
	for _, test := range []struct {
		text   string
		expect errData
	}{
		// Malformed syntax
		{"2020-01-01\n\t5h1", ErrorMalformedEntry().toErrData(2, 1, 3)},
		{"2020-01-01\n\tasdf Test 123", ErrorMalformedEntry().toErrData(2, 1, 4)},
		{"2020-01-01\n\t15:30", ErrorMalformedEntry().toErrData(2, 6, 1)},
		{"2020-01-01\n\t08:00-", ErrorMalformedEntry().toErrData(2, 7, 1)},
		{"2020-01-01\n\t08:00-asdf", ErrorMalformedEntry().toErrData(2, 7, 4)},
		{"2020-01-01\n\t08:00 - ?asdf", ErrorMalformedEntry().toErrData(2, 10, 4)},
		{"2020-01-01\n\t-18:00", ErrorMalformedEntry().toErrData(2, 1, 6)},
		{"2020-01-01\n\t15:30 Foo Bar Baz", ErrorMalformedEntry().toErrData(2, 7, 1)},

		// Logical errors
		{"2020-01-01\n\t08:00- ?\n\t09:00 - ?", ErrorDuplicateOpenRange().toErrData(3, 1, 9)},
		{"2020-01-01\n\t15:00 - 14:00", ErrorIllegalRange().toErrData(2, 1, 13)},
		{"2020-01-01\n\t12:76 - 13:00", ErrorMalformedEntry().toErrData(2, 1, 5)},
		{"2020-01-01\n\t12:00 - 44:00", ErrorMalformedEntry().toErrData(2, 9, 5)},
		{"2020-01-01\n\t23:00> - 25:61>", ErrorMalformedEntry().toErrData(2, 10, 6)},
		{"2020-01-01\n\t12:00> - 24:00>", ErrorMalformedEntry().toErrData(2, 10, 6)},
	} {
		rs, errs := Parse(test.text)
		require.Nil(t, rs, test.text)
		require.NotNil(t, errs, test.text)
		require.Len(t, errs, 1, test.text)
		assert.Equal(t, test.expect, toErrData(errs[0]), test.text)
	}
}
