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

func TestFindHashTagMatches(t *testing.T) {
	// `#` is stripped
	tags := NewTagSet("this", "#THAT", "numb3rs", "under_score")
	for _, x := range []struct {
		summary string
		matches []string
	}{
		{"#this at the beginning", []string{"this"}},
		{"#this, with punctuation afterwards", []string{"this"}},
		{"or at the end: #this", []string{"this"}},
		{"or #this in between", []string{"this"}},
		{"or all: #this and #that and #numb3rs", []string{"this", "that", "numb3rs"}},
		{"or #that as well (case-insensitive)", []string{"that"}},
		{"not case sensitive #THIS", []string{"this"}},
		{"can also contain #numb3rs!", []string{"numb3rs"}},
		{"or #under_score's", []string{"under_score"}},
		{"#some other tag", nil},
		{"#thisAndThat is similar but not the same", nil},
	} {
		matches := Summary(x.summary).MatchTags(tags)
		assert.Equal(t, matches, NewTagSet(x.matches...))
	}
}
