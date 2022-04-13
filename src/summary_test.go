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

func TestDetectsSummaryEquality(t *testing.T) {
	for _, x := range [][]string{
		nil,
		{""},
		{"a"},
		{"a", "b"},
	} {
		entrySummary1, _ := NewEntrySummary(x...)
		entrySummary2, _ := NewEntrySummary(x...)
		assert.True(t, entrySummary1.Equals(entrySummary2))
		assert.True(t, entrySummary2.Equals(entrySummary1))

		recordSummary1, _ := NewRecordSummary(x...)
		recordSummary2, _ := NewRecordSummary(x...)
		assert.True(t, recordSummary1.Equals(recordSummary2))
		assert.True(t, recordSummary2.Equals(recordSummary1))
	}
}

func TestEqualityOfEmptyEntrySummary(t *testing.T) {
	emptyEntrySummary, _ := NewEntrySummary()
	assert.True(t, emptyEntrySummary.Equals(nil))

	blankEntrySummary, _ := NewEntrySummary("")
	assert.True(t, blankEntrySummary.Equals(nil))
}

func TestDetectsSummaryInequality(t *testing.T) {
	for _, x := range []struct {
		ls1 []string
		ls2 []string
	}{
		{[]string{"a"}, nil},
		{[]string{"a"}, []string{"b"}},
		{[]string{"a"}, []string{"a", "b"}},
		{[]string{"a"}, []string{"a", ""}},
	} {
		{
			entrySummary1, _ := NewEntrySummary(x.ls1...)
			entrySummary2, _ := NewEntrySummary(x.ls2...)
			assert.False(t, entrySummary1.Equals(entrySummary2))
			assert.False(t, entrySummary2.Equals(entrySummary1))
		}
		{
			recordSummary1, _ := NewRecordSummary(x.ls1...)
			recordSummary2, _ := NewRecordSummary(x.ls2...)
			assert.False(t, recordSummary1.Equals(recordSummary2))
			assert.False(t, recordSummary2.Equals(recordSummary1))
		}
	}
}

func TestRecognisesAllTags(t *testing.T) {
	recordSummary, _ := NewRecordSummary(
		"Hello #world, I feel",
		"(super #GREAT) today #123_test: #234-foo!",
		"#太陽 #λουλούδι #पहाड #мир #Léift #ΓΕΙΑ-ΣΑΣ",
	)

	assert.Equal(t, recordSummary.Tags().ToStrings(), []string{
		"#123_test", "#234-foo", "#great", "#léift", "#world", "#γεια-σασ", "#λουλούδι", "#мир", "#पह", "#太陽",
	})

	assert.True(t, recordSummary.Tags().Contains(NewTagOrPanic("123_test", "")))
	assert.True(t, recordSummary.Tags().Contains(NewTagOrPanic("234-foo", "")))
	assert.True(t, recordSummary.Tags().Contains(NewTagOrPanic("太陽", "")))
	assert.True(t, recordSummary.Tags().Contains(NewTagOrPanic("λουλούδι", "")))
	assert.True(t, recordSummary.Tags().Contains(NewTagOrPanic("γεια-σασ", "")))
	assert.True(t, recordSummary.Tags().Contains(NewTagOrPanic("GREAT", "")))
	assert.True(t, recordSummary.Tags().Contains(NewTagOrPanic("Great", "")))
	assert.True(t, recordSummary.Tags().Contains(NewTagOrPanic("great", "")))
	assert.True(t, recordSummary.Tags().Contains(NewTagOrPanic("world", "")))

	assert.False(t, recordSummary.Tags().Contains(NewTagOrPanic("foo", "")))
	assert.False(t, recordSummary.Tags().Contains(NewTagOrPanic("test", "")))
	assert.False(t, recordSummary.Tags().Contains(NewTagOrPanic("test", "")))
	assert.False(t, recordSummary.Tags().Contains(NewTagOrPanic("ടെലിഫോണ്", "")))
	assert.False(t, recordSummary.Tags().Contains(NewTagOrPanic("123", "")))
	assert.False(t, recordSummary.Tags().Contains(NewTagOrPanic("wor", "")))
	assert.False(t, recordSummary.Tags().Contains(NewTagOrPanic("super", "")))
	assert.False(t, recordSummary.Tags().Contains(NewTagOrPanic("маркуч", "")))
	assert.False(t, recordSummary.Tags().Contains(NewTagOrPanic("grea", "")))
	assert.False(t, recordSummary.Tags().Contains(NewTagOrPanic("blabla", "")))

	entrySummary, _ := NewEntrySummary("Hello #world, I feel #great #TODAY")
	assert.Equal(t, entrySummary.Tags().ToStrings(), []string{"#great", "#today", "#world"})
}
