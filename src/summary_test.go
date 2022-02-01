package klog

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestCreatesEmptySummary(t *testing.T) {
	recordSummary, rErr := NewRecordSummary()
	require.Nil(t, rErr)
	assert.Nil(t, recordSummary.Lines())
	assert.Empty(t, recordSummary.Tags())

	entrySummary, eErr := NewEntrySummary()
	require.Nil(t, eErr)
	assert.Nil(t, entrySummary.Lines())
	assert.Empty(t, entrySummary.Tags())
}

func TestCreatesValidSingleLineSummary(t *testing.T) {
	recordSummary, rErr := NewRecordSummary("First line")
	require.Nil(t, rErr)
	assert.Equal(t, []string{"First line"}, recordSummary.Lines())
	assert.Empty(t, recordSummary.Tags())

	entrySummary, eErr := NewEntrySummary("First line")
	require.Nil(t, eErr)
	assert.Equal(t, []string{"First line"}, entrySummary.Lines())
	assert.Empty(t, entrySummary.Tags())
}

func TestCreatesValidMultilineSummary(t *testing.T) {
	recordSummary, rErr := NewRecordSummary("First line", "Second line")
	require.Nil(t, rErr)
	assert.Equal(t, []string{"First line", "Second line"}, recordSummary.Lines())
	assert.Empty(t, recordSummary.Tags())

	entrySummary, eErr := NewEntrySummary("First line", "Second line")
	require.Nil(t, eErr)
	assert.Equal(t, []string{"First line", "Second line"}, entrySummary.Lines())
	assert.Empty(t, entrySummary.Tags())
}

func TestRecordSummaryCannotContainBlankLines(t *testing.T) {
	for _, l := range [][]string{
		{""},
		{"     "},
		{"\u00a0\u00a0\u00a0\u00a0"},
		{"Foo", "\u00a0\u00a0\u00a0\u00a0"},
		{"\t\t"},
		{"Hello", "     ", "Foo"},
		{"Hello", "\t", "Foo"},
		{"Hello", "", "Foo"},
		{"Hello", "Foo", ""},
	} {
		recordSummary, err := NewRecordSummary(l...)
		require.Error(t, err)
		require.Nil(t, recordSummary)
	}
}

func TestRecordSummaryCannotContainWhitespaceAtBeginningOfLine(t *testing.T) {
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

func TestEntrySummaryCanStartWithBlankOrEmptyLine(t *testing.T) {
	for _, l := range [][]string{
		{"", "Foo"},
		{" ", "Foo", "Bar"},
		{"\t", " Foo"},
		{"\u00a0", "\tFoo     "},
		{"\u00a0\t     \t ", "   Foo", "\u00a0Baz \t"},
	} {
		entrySummary, err := NewEntrySummary(l...)
		require.Nil(t, err)
		require.NotNil(t, entrySummary)
	}
}

func TestEntrySummaryCannotContainSubsequentBlankLines(t *testing.T) {
	for _, l := range [][]string{
		{"Foo", ""},
		{"Foo", "     "},
		{"Foo", "\u00a0\u00a0\u00a0\u00a0"},
		{"Foo", "\t\t"},
		{"Hello", "     ", "Foo"},
		{"Hello", "\t", "Foo"},
		{"Hello", "", "Foo"},
		{"Hello", "Foo", ""},
	} {
		entrySummary, err := NewEntrySummary(l...)
		require.Error(t, err)
		require.Nil(t, entrySummary)
	}
}

func TestRecognisesAllTags(t *testing.T) {
	recordSummary, _ := NewRecordSummary("Hello #world, I feel", "#GREAT-ish today #123_test!")
	assert.Equal(t, recordSummary.Tags().ToStrings(), []string{"#123_test", "#great", "#world"})
	assert.True(t, recordSummary.Tags().Contains("#123_test"))
	assert.True(t, recordSummary.Tags().Contains("great"))
	assert.True(t, recordSummary.Tags().Contains("world"))

	entrySummary, _ := NewEntrySummary("Hello #world, I feel #great #TODAY")
	assert.Equal(t, entrySummary.Tags().ToStrings(), []string{"#great", "#today", "#world"})
}

func TestPerformsFuzzyMatching(t *testing.T) {
	summary, _ := NewRecordSummary("Hello #world, I feel #GREAT-ish today #123_test!")
	assert.True(t, summary.Tags().Contains("#123_..."))
	assert.True(t, summary.Tags().Contains("GR..."))
	assert.True(t, summary.Tags().Contains("WoRl..."))
	assert.False(t, summary.Tags().Contains("worl"))
}
