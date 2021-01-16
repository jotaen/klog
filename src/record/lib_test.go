package record

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestSumUpTimes(t *testing.T) {
	r := NewRecord(Ɀ_Date_(2020, 1, 1))
	r.AddDuration(NewDuration(1, 0), "")
	r.AddDuration(NewDuration(2, 0), "")
	assert.Equal(t, NewDuration(3, 0), Total(r))
}

func TestSumUpZeroIfNoTimesAvailable(t *testing.T) {
	r := NewRecord(Ɀ_Date_(2020, 1, 1))
	assert.Equal(t, NewDuration(0, 0), Total(r))
}

func TestSumUpRanges(t *testing.T) {
	range1 := Ɀ_Range_(Ɀ_Time_(9, 7), Ɀ_Time_(12, 59))
	range2 := Ɀ_Range_(Ɀ_Time_(13, 49), Ɀ_Time_(17, 12))
	r := NewRecord(Ɀ_Date_(2020, 1, 1))
	r.AddRange(range1, "")
	r.AddRange(range2, "")
	assert.Equal(t, NewDuration(7, 15), Total(r))
}

func TestSumUpTimesAndRanges(t *testing.T) {
	range1 := Ɀ_Range_(Ɀ_Time_(8, 0), Ɀ_Time_(12, 0))
	r := NewRecord(Ɀ_Date_(2020, 1, 1))
	r.AddDuration(NewDuration(1, 33), "")
	r.AddRange(range1, "")
	assert.Equal(t, NewDuration(5, 33), Total(r))
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
	r := NewRecord(Ɀ_Date_(2020, 1, 1))
	_ = r.SetSummary("This and #that, and other stuff as well")
	r.AddDuration(NewDuration(2, 0), "Foo")
	r.AddRange(Ɀ_Range_(Ɀ_Time_(13, 49), Ɀ_Time_(17, 12)), "Bar")
	es := FindEntriesWithHashtags(TagList("that"), r)
	require.Len(t, es, 2)
	assert.Equal(t, es[0].SummaryAsString(), "Foo")
	assert.Equal(t, es[1].SummaryAsString(), "Bar")
}

func TestHashTagReturnsEntriesThatMatch(t *testing.T) {
	r := NewRecord(Ɀ_Date_(2020, 1, 1))
	_ = r.SetSummary("This and #that, and other stuff as well")
	r.AddDuration(NewDuration(2, 0), "Foo #fizzbuzz")
	r.AddRange(Ɀ_Range_(Ɀ_Time_(13, 49), Ɀ_Time_(17, 12)), "Bar #barbaz")
	es := FindEntriesWithHashtags(TagList("fizzbuzz"), r)
	require.Len(t, es, 1)
	assert.Equal(t, es[0].SummaryAsString(), "Foo #fizzbuzz")
}
