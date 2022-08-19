package parser

import (
	"github.com/jotaen/klog/klog"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestParseMinimalDocument(t *testing.T) {
	text := `2000-01-01`
	rs, errs := Parse(text)
	require.Nil(t, errs)
	require.Len(t, rs, 1)
	assert.Equal(t, klog.Ɀ_Date_(2000, 1, 1), rs[0].Date())
}

func TestParseMultipleRecords(t *testing.T) {
	text := `
1999-05-31

1999-06-03
  1h
`
	rs, errs := Parse(text)
	require.Nil(t, errs)
	require.Len(t, rs, 2)

	assert.Equal(t, klog.Ɀ_Date_(1999, 5, 31), rs[0].Date())
	assert.Len(t, rs[0].Entries(), 0)

	assert.Equal(t, klog.Ɀ_Date_(1999, 6, 3), rs[1].Date())
	assert.Len(t, rs[1].Entries(), 1)
}

func TestParseCompleteRecord(t *testing.T) {
	text := `
1970-08-29 (8h15m!)
Record summary with
multiple lines of text
    1h
    1h1m Duration with summary
    1h2m Duration with
        multiline summary
    8:00-9:30
    9:00-10:31 Range with summary
    10:00-11:32 Range with multiple
        lines of
          summary text
    11:00-?
        Open range
`
	rs, errs := Parse(text)
	require.Nil(t, errs)
	require.Len(t, rs, 1)

	r := rs[0]
	assert.Equal(t, klog.Ɀ_Date_(1970, 8, 29), r.Date())
	assert.Equal(t, klog.Ɀ_RecordSummary_("Record summary with", "multiple lines of text"), r.Summary())
	assert.Equal(t, klog.NewDuration(8, 15).InMinutes(), r.ShouldTotal().InMinutes())

	assert.Len(t, r.Entries(), 7)

	assert.Equal(t, klog.NewDuration(1, 0).InMinutes(), r.Entries()[0].Duration().InMinutes())
	assert.Equal(t, klog.Ɀ_EntrySummary_(""), r.Entries()[0].Summary())

	assert.Equal(t, klog.NewDuration(1, 1).InMinutes(), r.Entries()[1].Duration().InMinutes())
	assert.Equal(t, klog.Ɀ_EntrySummary_("Duration with summary"), r.Entries()[1].Summary())

	assert.Equal(t, klog.NewDuration(1, 2).InMinutes(), r.Entries()[2].Duration().InMinutes())
	assert.Equal(t, klog.Ɀ_EntrySummary_("Duration with", "multiline summary"), r.Entries()[2].Summary())

	assert.Equal(t, klog.NewDuration(1, 30).InMinutes(), r.Entries()[3].Duration().InMinutes())
	assert.Equal(t, klog.Ɀ_EntrySummary_(""), r.Entries()[3].Summary())

	assert.Equal(t, klog.NewDuration(1, 31).InMinutes(), r.Entries()[4].Duration().InMinutes())
	assert.Equal(t, klog.Ɀ_EntrySummary_("Range with summary"), r.Entries()[4].Summary())

	assert.Equal(t, klog.NewDuration(1, 32).InMinutes(), r.Entries()[5].Duration().InMinutes())
	assert.Equal(t, klog.Ɀ_EntrySummary_("Range with multiple", "lines of", "  summary text"), r.Entries()[5].Summary())

	assert.Equal(t, klog.NewDuration(0, 0).InMinutes(), r.Entries()[6].Duration().InMinutes())
	assert.Equal(t, klog.Ɀ_EntrySummary_("", "Open range"), r.Entries()[6].Summary())
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
	text := "2018-01-01\n    1h\n" + "    \n" + "2019-01-01\n     \n2019-01-02"
	rs, errs := Parse(text)
	require.Nil(t, errs)
	require.Len(t, rs, 3)
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
	for _, x := range []string{
		"2000-01-01\n\t8h",
		"2000-01-01\n\t8h\n\t15m",
		"2000-05-31\n  6h",
		"2000-05-31\n  6h\n  20m",
		"2000-05-31\n   6h",
		"2000-05-31\n    6h",
	} {
		rs, errs := Parse(x)
		require.Nil(t, errs)
		require.Len(t, rs, 1)
	}
}

func TestParseDocumentSucceedsWithCorrectEntryValues(t *testing.T) {
	for _, test := range []struct {
		text        string
		expectEntry any
	}{
		// Durations
		{"1234-12-12\n\t5h", klog.NewDuration(5, 0)},
		{"1234-12-12\n\t2m", klog.NewDuration(0, 2)},
		{"1234-12-12\n\t2h30m", klog.NewDuration(2, 30)},

		// Durations with sign
		{"1234-12-12\n\t+5h", klog.NewDuration(5, 0)},
		{"1234-12-12\n\t+2h30m", klog.NewDuration(2, 30)},
		{"1234-12-12\n\t+2m", klog.NewDuration(0, 2)},
		{"1234-12-12\n\t-5h", klog.NewDuration(-5, -0)},
		{"1234-12-12\n\t-2h30m", klog.NewDuration(-2, -30)},
		{"1234-12-12\n\t-2m", klog.NewDuration(-0, -2)},

		// Ranges
		{"1234-12-12\n\t3:05 - 11:59", klog.Ɀ_Range_(klog.Ɀ_Time_(3, 5), klog.Ɀ_Time_(11, 59))},
		{"1234-12-12\n\t22:00 - 24:00", klog.Ɀ_Range_(klog.Ɀ_Time_(22, 0), klog.Ɀ_TimeTomorrow_(0, 0))},
		{"1234-12-12\n\t9:00am - 1:43pm", klog.Ɀ_Range_(klog.Ɀ_IsAmPm_(klog.Ɀ_Time_(9, 00)), klog.Ɀ_IsAmPm_(klog.Ɀ_Time_(13, 43)))},
		{"1234-12-12\n\t9:00am-1:43pm", klog.Ɀ_Range_(klog.Ɀ_IsAmPm_(klog.Ɀ_Time_(9, 00)), klog.Ɀ_IsAmPm_(klog.Ɀ_Time_(13, 43)))},
		{"1234-12-12\n\t9:00am-9:05", klog.Ɀ_Range_(klog.Ɀ_IsAmPm_(klog.Ɀ_Time_(9, 00)), klog.Ɀ_Time_(9, 05))},

		// Ranges with shifted times
		{"1234-12-12\n\t9:00am-8:12am>", klog.Ɀ_Range_(klog.Ɀ_IsAmPm_(klog.Ɀ_Time_(9, 00)), klog.Ɀ_IsAmPm_(klog.Ɀ_TimeTomorrow_(8, 12)))},
		{"1234-12-12\n\t<22:00 - <24:00", klog.Ɀ_Range_(klog.Ɀ_TimeYesterday_(22, 0), klog.Ɀ_Time_(0, 0))},
		{"1234-12-12\n\t<23:30 - 0:10", klog.Ɀ_Range_(klog.Ɀ_TimeYesterday_(23, 30), klog.Ɀ_Time_(0, 10))},
		{"1234-12-12\n\t22:17 - 1:00>", klog.Ɀ_Range_(klog.Ɀ_Time_(22, 17), klog.Ɀ_TimeTomorrow_(1, 00))},
		{"1234-12-12\n\t<23:00-1:00>", klog.Ɀ_Range_(klog.Ɀ_TimeYesterday_(23, 00), klog.Ɀ_TimeTomorrow_(1, 00))},
		{"1234-12-12\n\t<23:00-<23:10", klog.Ɀ_Range_(klog.Ɀ_TimeYesterday_(23, 00), klog.Ɀ_TimeYesterday_(23, 10))},
		{"1234-12-12\n\t12:01>-13:59>", klog.Ɀ_Range_(klog.Ɀ_TimeTomorrow_(12, 01), klog.Ɀ_TimeTomorrow_(13, 59))},

		// Open ranges
		{"1234-12-12\n\t12:01 - ?", klog.NewOpenRange(klog.Ɀ_Time_(12, 1))},
		{"1234-12-12\n\t6:45pm - ?", klog.NewOpenRange(klog.Ɀ_IsAmPm_(klog.Ɀ_Time_(18, 45)))},
		{"1234-12-12\n\t18:45 - ???", klog.NewOpenRange(klog.Ɀ_Time_(18, 45))},
		{"1234-12-12\n\t<3:12-??????", klog.NewOpenRange(klog.Ɀ_TimeYesterday_(3, 12))},
	} {
		rs, errs := Parse(test.text)
		require.Nil(t, errs, test.text)
		require.Len(t, rs, 1, test.text)
		require.Len(t, rs[0].Entries(), 1, test.text)
		value := klog.Unbox(&rs[0].Entries()[0],
			func(r klog.Range) any { return r },
			func(d klog.Duration) any { return d },
			func(o klog.OpenRange) any { return o },
		)
		assert.Equal(t, test.expectEntry, value, test.text)
	}
}

func TestParsesDocumentsWithEntrySummaries(t *testing.T) {
	for _, test := range []struct {
		text          string
		expectEntry   any
		expectSummary klog.EntrySummary
	}{
		// Single line entries
		{"1234-12-12\n\t5h Some remark", klog.NewDuration(5, 0), klog.Ɀ_EntrySummary_("Some remark")},
		{"1234-12-12\n\t3:05 - 11:59 Did this and that", klog.Ɀ_Range_(klog.Ɀ_Time_(3, 5), klog.Ɀ_Time_(11, 59)), klog.Ɀ_EntrySummary_("Did this and that")},
		{"1234-12-12\n\t9:00am-8:12am> Things", klog.Ɀ_Range_(klog.Ɀ_IsAmPm_(klog.Ɀ_Time_(9, 00)), klog.Ɀ_IsAmPm_(klog.Ɀ_TimeTomorrow_(8, 12))), klog.Ɀ_EntrySummary_("Things")},
		{"1234-12-12\n\t18:45 - ? Just started something", klog.NewOpenRange(klog.Ɀ_Time_(18, 45)), klog.Ɀ_EntrySummary_("Just started something")},
		{"1234-12-12\n\t5h    Some remark", klog.NewDuration(5, 0), klog.Ɀ_EntrySummary_("   Some remark")},
		{"1234-12-12\n\t5h\tSome remark", klog.NewDuration(5, 0), klog.Ɀ_EntrySummary_("Some remark")},
		{"1234-12-12\n\t9:00am-9:05 Mixed styles", klog.Ɀ_Range_(klog.Ɀ_IsAmPm_(klog.Ɀ_Time_(9, 00)), klog.Ɀ_Time_(9, 05)), klog.Ɀ_EntrySummary_("Mixed styles")},
		{"1234-12-12\n\t3:05 - 11:59\tFoo", klog.Ɀ_Range_(klog.Ɀ_Time_(3, 5), klog.Ɀ_Time_(11, 59)), klog.Ɀ_EntrySummary_("Foo")},
		{"1234-12-12\n\t<22:00 - <24:00\tFoo", klog.Ɀ_Range_(klog.Ɀ_TimeYesterday_(22, 0), klog.Ɀ_Time_(0, 0)), klog.Ɀ_EntrySummary_("Foo")},
		{"1234-12-12\n\t22:00 - 24:00\tFoo", klog.Ɀ_Range_(klog.Ɀ_Time_(22, 0), klog.Ɀ_TimeTomorrow_(0, 0)), klog.Ɀ_EntrySummary_("Foo")},
		{"1234-12-12\n\t18:45 - ???       ASDF", klog.NewOpenRange(klog.Ɀ_Time_(18, 45)), klog.Ɀ_EntrySummary_("      ASDF")},
		{"1234-12-12\n\t18:45 - ?\tFoo", klog.NewOpenRange(klog.Ɀ_Time_(18, 45)), klog.Ɀ_EntrySummary_("Foo")},

		// Multiline-summary entries
		{"1234-12-12\n\t5h Some remark\n\t\twith more text", klog.NewDuration(5, 0), klog.Ɀ_EntrySummary_("Some remark", "with more text")},
		{"1234-12-12\n\t8:00-9:00 Some remark\n\t\twith more text", klog.Ɀ_Range_(klog.Ɀ_Time_(8, 00), klog.Ɀ_Time_(9, 00)), klog.Ɀ_EntrySummary_("Some remark", "with more text")},
		{"1234-12-12\n\t5h Some remark\n\t\twith\n\t\tmore\n\t\ttext", klog.NewDuration(5, 0), klog.Ɀ_EntrySummary_("Some remark", "with", "more", "text")},
		{"1234-12-12\n  5h Some remark\n    with more text", klog.NewDuration(5, 0), klog.Ɀ_EntrySummary_("Some remark", "with more text")},
		{"1234-12-12\n   5h Some remark\n      with more text", klog.NewDuration(5, 0), klog.Ɀ_EntrySummary_("Some remark", "with more text")},
		{"1234-12-12\n    5h Some remark\n        with more text", klog.NewDuration(5, 0), klog.Ɀ_EntrySummary_("Some remark", "with more text")},
		{"1234-12-12\n    5h Some remark\n        with\n        more\n        text", klog.NewDuration(5, 0), klog.Ɀ_EntrySummary_("Some remark", "with", "more", "text")},

		// Multiline-summary entries where first summary line is empty
		{"1234-12-12\n\t5h\n\t\twith more text", klog.NewDuration(5, 0), klog.Ɀ_EntrySummary_("", "with more text")},
		{"1234-12-12\n\t5h \n\t\twith more text", klog.NewDuration(5, 0), klog.Ɀ_EntrySummary_("", "with more text")},
		{"1234-12-12\n\t5h  \n\t\twith more text", klog.NewDuration(5, 0), klog.Ɀ_EntrySummary_(" ", "with more text")},
		{"1234-12-12\n\t5h\n\t\t with more text", klog.NewDuration(5, 0), klog.Ɀ_EntrySummary_("", " with more text")},
		{"1234-12-12\n\t5h\n\t\t\twith more text", klog.NewDuration(5, 0), klog.Ɀ_EntrySummary_("", "\twith more text")},
	} {
		rs, errs := Parse(test.text)
		require.Nil(t, errs, test.text)
		require.Len(t, rs, 1, test.text)
		require.Len(t, rs[0].Entries(), 1, test.text)
		value := klog.Unbox(&rs[0].Entries()[0],
			func(r klog.Range) any { return r },
			func(d klog.Duration) any { return d },
			func(o klog.OpenRange) any { return o },
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

		// Mixed styles for entries within one record:
		{"2020-01-01\n    8h\n\t2h", ErrorIllegalIndentation().toErrData(3, 0, 3)},
		{"2020-01-01\n  8h\n\t2h", ErrorIllegalIndentation().toErrData(3, 0, 3)},
		{"2020-01-01\n\t8h\n    2h", ErrorIllegalIndentation().toErrData(3, 0, 6)},
		{"2020-01-01\n\t8h\n  2h", ErrorIllegalIndentation().toErrData(3, 0, 4)},
		{"2020-01-01\n    8h\n  2h", ErrorIllegalIndentation().toErrData(3, 0, 4)},
		{"2020-01-01\n  8h\n   2h", ErrorIllegalIndentation().toErrData(3, 0, 5)},

		// Mixed styles for entry summaries within one record:
		{"2020-01-01\n  8h Foo\n\tbar baz", ErrorIllegalIndentation().toErrData(3, 0, 8)},
		{"2020-01-01\n    8h Foo\n       bar baz", ErrorIllegalIndentation().toErrData(3, 0, 14)},
		{"2020-01-01\n    8h Foo\n      bar baz", ErrorIllegalIndentation().toErrData(3, 0, 13)},
		{"2020-01-01\n    8h Foo\n    \tbar baz", ErrorIllegalIndentation().toErrData(3, 0, 12)},
		{"2020-01-01\n   8h Foo\n     bar baz", ErrorIllegalIndentation().toErrData(3, 0, 12)},
		{"2020-01-01\n  8h Foo\n   bar baz", ErrorIllegalIndentation().toErrData(3, 0, 10)},
		{"2020-01-01\n  8h\n  8h Foo\n   bar baz", ErrorIllegalIndentation().toErrData(4, 0, 10)},
	} {
		rs, errs := Parse(test.text)
		require.Nil(t, rs, test.text)
		require.NotNil(t, errs, test.text)
		require.Len(t, errs, 1, test.text)
		assert.Equal(t, test.expect, toErrData(errs[0]), test.text)
	}
}

func TestAcceptMixingIndentationStylesAcrossDifferentRecords(t *testing.T) {
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
