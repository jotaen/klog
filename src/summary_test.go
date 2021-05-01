package klog

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

func TestRecognisesAllTags(t *testing.T) {
	s := Summary("Hello #world, I feel #GREAT-ish today #123_test!")
	assert.Equal(t, s.Tags().ToStrings(), []string{"#123_test", "#great", "#world"})
	assert.True(t, s.Tags().Contains("#123_test"))
	assert.True(t, s.Tags().Contains("great"))
	assert.True(t, s.Tags().Contains("world"))
}

func TestPerformsFuzzyMatching(t *testing.T) {
	s := Summary("Hello #world, I feel #GREAT-ish today #123_test!")
	assert.True(t, s.Tags().Contains("#123_..."))
	assert.True(t, s.Tags().Contains("GR..."))
	assert.True(t, s.Tags().Contains("WoRl..."))
	assert.False(t, s.Tags().Contains("worl"))
}
