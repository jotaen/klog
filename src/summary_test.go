package klog

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestCreatesEmptySummary(t *testing.T) {
	summary, err := NewRecordSummary()
	require.Nil(t, err)
	assert.Nil(t, summary.Lines())
	assert.True(t, summary.IsEmpty())
	assert.Empty(t, summary.Tags())
}

func TestCreatesValidMultilineSummary(t *testing.T) {
	summary, err := NewRecordSummary("First line", "Second line")
	require.Nil(t, err)
	assert.Equal(t, []string{"First line", "Second line"}, summary.Lines())
	assert.False(t, summary.IsEmpty())
	assert.Empty(t, summary.Tags())
}

}

func TestSummaryCannotContainWhitespaceAtBeginningOfLine(t *testing.T) {
	r := NewRecord(Ɀ_Date_(2020, 1, 1))
	require.Error(t, r.SetSummary(NewRecordSummary("Hello", " World")))
	require.Error(t, r.SetSummary(NewRecordSummary(" Hello")))
	assert.Equal(t, NewRecordSummary(), r.Summary()) // Still empty
}

func TestRecognisesAllTags(t *testing.T) {
	summary, _ := NewRecordSummary("Hello #world, I feel", "#GREAT-ish today #123_test!")
	assert.Equal(t, summary.Tags().ToStrings(), []string{"#123_test", "#great", "#world"})
	assert.True(t, summary.Tags().Contains("#123_test"))
	assert.True(t, summary.Tags().Contains("great"))
	assert.True(t, summary.Tags().Contains("world"))
}

func TestPerformsFuzzyMatching(t *testing.T) {
	summary, _ := NewRecordSummary("Hello #world, I feel #GREAT-ish today #123_test!")
	assert.True(t, summary.Tags().Contains("#123_..."))
	assert.True(t, summary.Tags().Contains("GR..."))
	assert.True(t, summary.Tags().Contains("WoRl..."))
	assert.False(t, summary.Tags().Contains("worl"))
}
