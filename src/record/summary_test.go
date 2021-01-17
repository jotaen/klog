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
