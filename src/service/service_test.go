package service

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"klog/record"
	"testing"
)

func TestSumUpTimes(t *testing.T) {
	r := record.NewRecord(record.Ɀ_Date_(2020, 1, 1))
	r.AddDuration(record.NewDuration(1, 0), "")
	r.AddDuration(record.NewDuration(2, 0), "")
	assert.Equal(t, record.NewDuration(3, 0), Total(r))
}

func TestSumUpZeroIfNoTimesAvailable(t *testing.T) {
	r := record.NewRecord(record.Ɀ_Date_(2020, 1, 1))
	assert.Equal(t, record.NewDuration(0, 0), Total(r))
}

func TestSumUpRanges(t *testing.T) {
	range1 := record.Ɀ_Range_(record.Ɀ_Time_(9, 7), record.Ɀ_Time_(12, 59))
	range2 := record.Ɀ_Range_(record.Ɀ_Time_(13, 49), record.Ɀ_Time_(17, 12))
	r := record.NewRecord(record.Ɀ_Date_(2020, 1, 1))
	r.AddRange(range1, "")
	r.AddRange(range2, "")
	assert.Equal(t, record.NewDuration(7, 15), Total(r))
}

func TestSumUpTimesAndRanges(t *testing.T) {
	range1 := record.Ɀ_Range_(record.Ɀ_Time_(8, 0), record.Ɀ_Time_(12, 0))
	r := record.NewRecord(record.Ɀ_Date_(2020, 1, 1))
	r.AddDuration(record.NewDuration(1, 33), "")
	r.AddRange(range1, "")
	assert.Equal(t, record.NewDuration(5, 33), Total(r))
}

func TestHashTagMatches(t *testing.T) {
	tags := TagList("this", "THAT")
	for _, txt := range []string{
		"#this at the beginning",
		"#this, with punctuation afterwards",
		"or at the end: #this",
		"or #this in between",
		"or both #this and #that",
		"or #that as well (case-insensitive)",
		"not case sensitive #THIS",
	} {
		isMatch := ContainsOneOfTags(tags, txt)
		assert.True(t, isMatch)
	}
}

func TestHashTagDoesNotMatch(t *testing.T) {
	tags := TagList("this", "that")
	for _, txt := range []string{
		"#some other tag",
		"#thisAndThat is not the same",
	} {
		isMatch := ContainsOneOfTags(tags, txt)
		assert.False(t, isMatch)
	}
}

func TestHashTagAllEntriesAreReturnedIfMatchIsInSummary(t *testing.T) {
	r := record.NewRecord(record.Ɀ_Date_(2020, 1, 1))
	_ = r.SetSummary("This and #that, and other stuff as well")
	r.AddDuration(record.NewDuration(2, 0), "Foo")
	r.AddRange(record.Ɀ_Range_(record.Ɀ_Time_(13, 49), record.Ɀ_Time_(17, 12)), "Bar")
	es := FindEntriesWithHashtags(TagList("that"), r)
	require.Len(t, es, 2)
	assert.Equal(t, record.Summary("Foo"), es[0].Summary())
	assert.Equal(t, record.Summary("Bar"), es[1].Summary())
}

func TestHashTagReturnsEntriesThatMatch(t *testing.T) {
	r := record.NewRecord(record.Ɀ_Date_(2020, 1, 1))
	_ = r.SetSummary("This and #that, and other stuff as well")
	r.AddDuration(record.NewDuration(2, 0), "Foo #fizzbuzz")
	r.AddRange(record.Ɀ_Range_(record.Ɀ_Time_(13, 49), record.Ɀ_Time_(17, 12)), "Bar #barbaz")
	es := FindEntriesWithHashtags(TagList("fizzbuzz"), r)
	require.Len(t, es, 1)
	assert.Equal(t, record.Summary("Foo #fizzbuzz"), es[0].Summary())
}
