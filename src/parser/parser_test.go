package parser

import (
	"github.com/jotaen/klog/src"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestParseMinimalDocument(t *testing.T) {
	text := `2000-01-01`
	pr, errs := Parse(text)
	require.Nil(t, errs)
	require.Len(t, pr.Records, 1)
	assert.Equal(t, klog.Ɀ_Date_(2000, 1, 1), pr.Records[0].Date())
}

func TestParseMultipleRecords(t *testing.T) {
	text := `
1999-05-31

1999-06-03 (8h15m!)
Empty
`
	pr, errs := Parse(text)
	require.Nil(t, errs)
	require.Len(t, pr.Records, 2)

	assert.Equal(t, klog.Ɀ_Date_(1999, 5, 31), pr.Records[0].Date())
	assert.Equal(t, klog.NewSummary(), pr.Records[0].Summary())
	assert.Len(t, pr.Records[0].Entries(), 0)

	assert.Equal(t, klog.Ɀ_Date_(1999, 6, 3), pr.Records[1].Date())
	assert.Equal(t, klog.NewSummary("Empty"), pr.Records[1].Summary())
	assert.Equal(t, klog.NewDuration(8, 15).InMinutes(), pr.Records[1].ShouldTotal().InMinutes())
	assert.Len(t, pr.Records[1].Entries(), 0)
}

func TestParseEmptyOrBlankDocument(t *testing.T) {
	for _, text := range []string{
		"",
		"    ",
		"\n\n\n\n\n",
		"\n\t     \n \n         ",
	} {
		pr, errs := Parse(text)
		require.Nil(t, errs)
		require.Len(t, pr.Records, 0)
	}
}

func TestParseWindowsAndUnixLineEndings(t *testing.T) {
	text := "2000-01-01\r\n\r\n2000-01-02\n\n2000-01-03"
	pr, errs := Parse(text)
	require.Nil(t, errs)
	require.Len(t, pr.Records, 3)
}

func TestParseMultipleRecordsWhenBlankLineContainsWhitespace(t *testing.T) {
	text := "2018-01-01\n    1h\n" + "    \n" + "2019-01-01\n"
	pr, errs := Parse(text)
	require.Nil(t, errs)
	require.Len(t, pr.Records, 2)
}

func TestParseAlternativeFormatting(t *testing.T) {
	text := `
1999/05/31
	8:00-13:00

1999-05-31
	8:00am-1:00pm
`
	pr, errs := Parse(text)
	require.Nil(t, errs)
	require.Len(t, pr.Records, 2)

	assert.True(t, pr.Records[0].Date().IsEqualTo(pr.Records[1].Date()))
	assert.Equal(t, pr.Records[0].Entries()[0].Duration(), pr.Records[1].Entries()[0].Duration())
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
		pr, errs := Parse(test.text)
		require.Nil(t, errs)
		require.Len(t, pr.Records, 1)
	}
}

func TestParseDocumentSucceedsWithCorrectEntries(t *testing.T) {
	for _, test := range []struct {
		text          string
		expectEntry   interface{}
		expectSummary klog.Summary
	}{
		{"1234-12-12\n\t5h Some remark", klog.NewDuration(5, 0), klog.NewSummary("Some remark")},
		{"1234-12-12\n\t2h30m", klog.NewDuration(2, 30), klog.NewSummary()},
		{"1234-12-12\n\t2m", klog.NewDuration(0, 2), klog.NewSummary()},
		{"1234-12-12\n\t+5h", klog.NewDuration(5, 0), klog.NewSummary()},
		{"1234-12-12\n\t+2h30m", klog.NewDuration(2, 30), klog.NewSummary()},
		{"1234-12-12\n\t+2m", klog.NewDuration(0, 2), klog.NewSummary()},
		{"1234-12-12\n\t-5h", klog.NewDuration(-5, -0), klog.NewSummary()},
		{"1234-12-12\n\t-2h30m", klog.NewDuration(-2, -30), klog.NewSummary()},
		{"1234-12-12\n\t-2m", klog.NewDuration(-0, -2), klog.NewSummary()},
		{"1234-12-12\n\t3:05 - 11:59 Did this and that", klog.Ɀ_Range_(klog.Ɀ_Time_(3, 5), klog.Ɀ_Time_(11, 59)), klog.NewSummary("Did this and that")},
		{"1234-12-12\n\t<23:30 - 0:10", klog.Ɀ_Range_(klog.Ɀ_TimeYesterday_(23, 30), klog.Ɀ_Time_(0, 10)), klog.NewSummary()},
		{"1234-12-12\n\t22:17 - 1:00>", klog.Ɀ_Range_(klog.Ɀ_Time_(22, 17), klog.Ɀ_TimeTomorrow_(1, 00)), klog.NewSummary()},
		{"1234-12-12\n\t18:45 - ? Just started something", klog.NewOpenRange(klog.Ɀ_Time_(18, 45)), klog.NewSummary("Just started something")},
		{"1234-12-12\n\t<3:12-??????", klog.NewOpenRange(klog.Ɀ_TimeYesterday_(3, 12)), klog.NewSummary()},
	} {
		pr, errs := Parse(test.text)
		require.Nil(t, errs, test.text)
		require.Len(t, pr.Records, 1, test.text)
		require.Len(t, pr.Records[0].Entries(), 1, test.text)
		value := pr.Records[0].Entries()[0].Unbox(
			func(r klog.Range) interface{} { return r },
			func(d klog.Duration) interface{} { return d },
			func(o klog.OpenRange) interface{} { return o },
		)
		assert.Equal(t, test.expectEntry, value, test.text)
		assert.Equal(t, test.expectSummary, pr.Records[0].Entries()[0].Summary(), test.text)
	}
}

func TestMalformedRecord(t *testing.T) {
	text := `
1999-05-31
	5h30m This and that
Why is there a summary at the end?
`
	pr, errs := Parse(text)
	require.Nil(t, pr)
	require.NotNil(t, errs)
	require.Len(t, errs.Get(), 1)
	assert.Equal(t, Err{id(ErrorIllegalIndentation), 4, 0, 34}, toErr(errs.Get()[0]))
}

func TestReportErrorsInHeadline(t *testing.T) {
	for _, test := range []struct {
		text   string
		expect Err
	}{
		{"Hello 123", Err{id(ErrorInvalidDate), 1, 0, 5}},
		{" 2020-01-01", Err{id(ErrorIllegalIndentation), 1, 0, 10}},
		{"   2020-01-01", Err{id(ErrorIllegalIndentation), 1, 0, 10}},
		{"2020-01-01 ()", Err{id(ErrorMalformedPropertiesSyntax), 1, 12, 1}},
		{"2020-01-01 (asdf)", Err{id(ErrorUnrecognisedProperty), 1, 12, 4}},
		{"2020-01-01 (asdf!)", Err{id(ErrorMalformedShouldTotal), 1, 12, 4}},
		{"2020-01-01 5h30m!", Err{id(ErrorUnrecognisedTextInHeadline), 1, 11, 6}},
		{"2020-01-01 (5h30m!", Err{id(ErrorMalformedPropertiesSyntax), 1, 18, 1}},
		{"2020-01-01 (", Err{id(ErrorMalformedPropertiesSyntax), 1, 12, 1}},
		{"2020-01-01 (5h!) foo", Err{id(ErrorUnrecognisedTextInHeadline), 1, 17, 3}},
		{"2020-01-01 (5h! asdf)", Err{id(ErrorUnrecognisedProperty), 1, 16, 4}},
		{"2020-01-01 (5h!!!)", Err{id(ErrorUnrecognisedProperty), 1, 15, 2}},
	} {
		pr, errs := Parse(test.text)
		require.Nil(t, pr)
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
	pr, errs := Parse(text)
	require.Nil(t, pr)
	require.NotNil(t, errs)
	require.Len(t, errs.Get(), 1)
	assert.Equal(t, Err{id(ErrorIllegalIndentation), 4, 0, 40}, toErr(errs.Get()[0]))
}

func TestReportErrorsInEntries(t *testing.T) {
	for _, test := range []struct {
		text   string
		expect Err
	}{
		{"2020-01-01\n\t5h1", Err{id(ErrorMalformedEntry), 2, 0, 3}},
		{"2020-01-01\n\tasdf Test 123", Err{id(ErrorMalformedEntry), 2, 0, 4}},
		{"2020-01-01\n\t15:30", Err{id(ErrorMalformedEntry), 2, 5, 1}},
		{"2020-01-01\n\t08:00-", Err{id(ErrorMalformedEntry), 2, 6, 1}},
		{"2020-01-01\n\t08:00-asdf", Err{id(ErrorMalformedEntry), 2, 6, 4}},
		{"2020-01-01\n\t08:00 - ?asdf", Err{id(ErrorMalformedEntry), 2, 9, 4}},
		{"2020-01-01\n\t08:00- ?\n\t09:00 - ?", Err{id(ErrorDuplicateOpenRange), 3, 0, 9}},
		{"2020-01-01\n\t15:00 - 14:00", Err{id(ErrorIllegalRange), 2, 0, 13}},
		{"2020-01-01\n\t-18:00", Err{id(ErrorMalformedEntry), 2, 0, 6}},
		{"2020-01-01\n\t15:30 Foo Bar Baz", Err{id(ErrorMalformedEntry), 2, 6, 1}},
		{"2020-01-01\n 8h", Err{id(ErrorIllegalIndentation), 2, 0, 2}},
		{"2020-01-01\n\t 8h", Err{id(ErrorIllegalIndentation), 2, 0, 2}},
		{"2020-01-01\n\t\t8h", Err{id(ErrorIllegalIndentation), 2, 0, 2}},
		{"2020-01-01\n     8h", Err{id(ErrorIllegalIndentation), 2, 0, 2}},
	} {
		pr, errs := Parse(test.text)
		require.Nil(t, pr, test.text)
		require.NotNil(t, errs, test.text)
		require.Len(t, errs.Get(), 1, test.text)
		assert.Equal(t, test.expect, toErr(errs.Get()[0]), test.text)
	}
}
