package record

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestSavesSummary(t *testing.T) {
	r := NewRecord(Ɀ_Date_(2020, 1, 1))
	err := r.SetSummary("Hello World")
	require.Nil(t, err)
	assert.Equal(t, Summary("Hello World"), r.Summary())
}

func TestSummaryCannotContainWhitespaceAtBeginningOfLine(t *testing.T) {
	r := NewRecord(Ɀ_Date_(2020, 1, 1))
	require.Error(t, r.SetSummary("Hello\n World"))
	require.Error(t, r.SetSummary(" Hello"))
	assert.Equal(t, Summary(""), r.Summary()) // Still empty
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
