package parser

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	. "klog/parser/engine"
	. "klog/record"
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
	assert.Equal(t, Ɀ_Date_(2000, 1, 1), rs[0].Date())
}

func TestAcceptTabOrSpacesAsIndentation(t *testing.T) {
	text := `
2000-01-01
	8h

2000-05-31
  6h

2000-05-31
   6h

2000-05-31
    6h
`
	rs, errs := Parse(text)
	require.Nil(t, errs)
	require.Len(t, rs, 4)
}

func TestParseMultipleRecords(t *testing.T) {
	text := `
1999-05-31

1999-06-03
Empty
`
	rs, errs := Parse(text)
	require.Nil(t, errs)
	require.Len(t, rs, 2)

	assert.Equal(t, rs[0].Date(), Ɀ_Date_(1999, 5, 31))
	assert.Equal(t, rs[0].Summary(), "")
	assert.Len(t, rs[0].Entries(), 0)

	assert.Equal(t, rs[1].Date(), Ɀ_Date_(1999, 6, 3))
	assert.Equal(t, "Empty", rs[1].Summary())
	assert.Len(t, rs[1].Entries(), 0)
}

func TestParseDocumentWithFullFeaturedEntry(t *testing.T) {
	text := `
2020-01-15 (5h30m!)
This is a
multiline summary
	5h Some remark
	2h30m
	2m
	-5h
	-2h30m
	-2m
	3:05 - 11:59 Did this and that
	<23:30 - 0:10
	22:17 - 1:00>
	18:45 - Just started something
`
	rs, errs := Parse(text)
	require.Nil(t, errs)
	require.Len(t, rs, 1)
	assert.Equal(t, Ɀ_Date_(2020, 1, 15), rs[0].Date())
	assert.Equal(t, NewDuration(5, 30), rs[0].ShouldTotal())
	assert.Equal(t, "This is a\nmultiline summary", rs[0].Summary())
	require.Len(t, rs[0].Entries(), 10)
	assert.Equal(t, NewDuration(5, 0), rs[0].Entries()[0].Value())
	assert.Equal(t, "Some remark", rs[0].Entries()[0].SummaryAsString())
	assert.Equal(t, NewDuration(2, 30), rs[0].Entries()[1].Value())
	assert.Equal(t, NewDuration(0, 2), rs[0].Entries()[2].Value())
	assert.Equal(t, NewDuration(-5, -0), rs[0].Entries()[3].Value())
	assert.Equal(t, NewDuration(-2, -30), rs[0].Entries()[4].Value())
	assert.Equal(t, NewDuration(-0, -2), rs[0].Entries()[5].Value())
	assert.Equal(t, Ɀ_Range_(Ɀ_Time_(3, 5), Ɀ_Time_(11, 59)), rs[0].Entries()[6].Value())
	assert.Equal(t, "Did this and that", rs[0].Entries()[6].SummaryAsString())
	assert.Equal(t, Ɀ_Range_(Ɀ_TimeYesterday_(23, 30), Ɀ_Time_(0, 10)), rs[0].Entries()[7].Value())
	assert.Equal(t, Ɀ_Range_(Ɀ_Time_(22, 17), Ɀ_TimeTomorrow_(1, 00)), rs[0].Entries()[8].Value())
	assert.Equal(t, Ɀ_Time_(18, 45), rs[0].Entries()[9].Value())
	assert.Equal(t, "Just started something", rs[0].Entries()[9].SummaryAsString())
}

func TestMalformedRecord(t *testing.T) {
	text := `
1999-05-31
	5h30m This and that
Why is there a summary at the end?
`
	rs, errs := Parse(text)
	require.Nil(t, rs)
	require.Len(t, errs, 1)
	assert.Equal(t, Err{ILLEGAL_INDENTATION, 4, 0, 34}, toErr(errs[0].(Error)))
}

func TestReportErrorsInHeadline(t *testing.T) {
	text := `Hello 123

 2020-01-01

2020-01-01 (asdf)

2020-01-01 5h30m!

2020-01-01 (5h30m!
`
	rs, errs := Parse(text)
	require.Nil(t, rs)
	require.Len(t, errs, 5)
	assert.Equal(t, Err{INVALID_VALUE, 1, 0, 5}, toErr(errs[0].(Error)))
	assert.Equal(t, Err{ILLEGAL_WHITESPACE, 3, 0, 1}, toErr(errs[1].(Error)))
	assert.Equal(t, Err{INVALID_VALUE, 5, 12, 5}, toErr(errs[2].(Error)))
	assert.Equal(t, Err{INVALID_VALUE, 7, 11, 6}, toErr(errs[3].(Error)))
	assert.Equal(t, Err{INVALID_VALUE, 9, 18, 1}, toErr(errs[4].(Error)))
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
	require.Len(t, errs, 1)
	assert.Equal(t, Err{INVALID_VALUE, 4, 0, 41}, toErr(errs[0].(Error)))
}

func TestReportErrorsInEntries(t *testing.T) {
	text := `
2020-01-01
	5h1

2020-01-01
	15:30

2020-01-01
	asdf Test 123

2020-01-01
	08:00-
	09:00 -

2020-01-01
	15:00 - 14:00
`
	rs, errs := Parse(text)
	require.Nil(t, rs)
	require.Len(t, errs, 5)
	assert.Equal(t, Err{INVALID_VALUE, 3, 0, 3}, toErr(errs[0].(Error)))
	assert.Equal(t, Err{INVALID_VALUE, 6, 5, 5}, toErr(errs[1].(Error)))
	assert.Equal(t, Err{INVALID_VALUE, 9, 0, 4}, toErr(errs[2].(Error)))
	assert.Equal(t, Err{DUPLICATE_OPEN_RANGE, 13, 0, 7}, toErr(errs[3].(Error)))
	assert.Equal(t, Err{ILLEGAL_RANGE, 16, 0, 13}, toErr(errs[4].(Error)))
}
