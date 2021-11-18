package klog

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestCreatesEmptySummary(t *testing.T) {
	recordSummary, err := NewRecordSummary()
	require.Nil(t, err)
	assert.Nil(t, recordSummary.Lines())
	assert.True(t, recordSummary.IsEmpty())
	assert.Empty(t, recordSummary.Tags())

	entrySummary := NewEntrySummary("")
	assert.Nil(t, entrySummary.Lines())
	assert.True(t, entrySummary.IsEmpty())
	assert.Empty(t, entrySummary.Tags())
}

func TestCreatesValidSummary(t *testing.T) {
	recordSummary, err := NewRecordSummary("First line", "Second line")
	require.Nil(t, err)
	assert.Equal(t, []string{"First line", "Second line"}, recordSummary.Lines())
	assert.False(t, recordSummary.IsEmpty())
	assert.Empty(t, recordSummary.Tags())

	entrySummary := NewEntrySummary("First line")
	assert.Equal(t, []string{"First line"}, entrySummary.Lines())
	assert.False(t, entrySummary.IsEmpty())
	assert.Empty(t, entrySummary.Tags())
}

func TestSummaryCannotContainBlankLines(t *testing.T) {
	for _, l := range [][]string{
		{""},
		{"     "},
		{"\u00a0\u00a0\u00a0\u00a0"},
		{"\t\t"},
		{"Hello", "     ", "Foo"},
		{"Hello", "", "Foo"},
		{"Hello", "Foo", ""},
	} {
		summary, err := NewRecordSummary(l...)
		require.Error(t, err)
		require.Nil(t, summary)
	}
}

func TestSummaryCannotContainWhitespaceAtBeginningOfLine(t *testing.T) {
	for _, l := range [][]string{
		{" Hello"},
		{"\u00a0Hello"},
		{"\u2000Hello"},
		{"\u2007Hello"},
		{"\tHello"},
		{"Hello", " World"},
		{"Hello", "\tWorld"},
		{"Hello", "\u00a0World"},
	} {
		summary, err := NewRecordSummary(l...)
		require.Error(t, err)
		require.Nil(t, summary)
	}
}

func TestRecognisesAllTags(t *testing.T) {
	recordSummary, _ := NewRecordSummary("Hello #world, I feel", "#GREAT-ish today #123_test!")
	assert.Equal(t, recordSummary.Tags().ToStrings(), []string{"#123_test", "#great", "#world"})
	assert.True(t, recordSummary.Tags().Contains("#123_test"))
	assert.True(t, recordSummary.Tags().Contains("great"))
	assert.True(t, recordSummary.Tags().Contains("world"))

	entrySummary := NewEntrySummary("Hello #world, I feel #great #TODAY")
	assert.Equal(t, entrySummary.Tags().ToStrings(), []string{"#great", "#today", "#world"})
}

func TestPerformsFuzzyMatching(t *testing.T) {
	summary, _ := NewRecordSummary("Hello #world, I feel #GREAT-ish today #123_test!")
	assert.True(t, summary.Tags().Contains("#123_..."))
	assert.True(t, summary.Tags().Contains("GR..."))
	assert.True(t, summary.Tags().Contains("WoRl..."))
	assert.False(t, summary.Tags().Contains("worl"))
}
