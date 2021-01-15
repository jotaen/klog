package parser

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	. "klog/record"
	"testing"
)

func TestParseEmptyDocument(t *testing.T) {
	text := ``
	rs, err := Parse(text)
	require.Nil(t, err)
	require.Nil(t, rs)
}

func TestParseBlankDocument(t *testing.T) {
	text := `
 
`
	rs, err := Parse(text)
	require.Nil(t, err)
	require.Nil(t, rs)
}

func TestParseMultipleRecords(t *testing.T) {
	text := `
1999-05-31


1999-06-03
Empty
`
	rs, err := Parse(text)
	require.Nil(t, err)
	require.NotNil(t, rs)
	require.Len(t, rs, 2)

	assert.Equal(t, rs[0].Date(), Ɀ_Date_(1999, 5, 31))
	assert.Equal(t, rs[0].Summary(), "")
	assert.Len(t, rs[0].Entries(), 0)

	assert.Equal(t, rs[1].Date(), Ɀ_Date_(1999, 6, 3))
	assert.Equal(t, "Empty", rs[1].Summary())
	assert.Len(t, rs[1].Entries(), 0)
}

func TestParseDocumentWithSingleEntry(t *testing.T) {
	text := `
2020-01-15 (5h30m!)
This is a
multiline summary
`
	rs, err := Parse(text)
	require.Nil(t, err)
	require.NotNil(t, rs)
	require.Len(t, rs, 1)
	assert.Equal(t, Ɀ_Date_(2020, 1, 15), rs[0].Date())
	assert.Equal(t, "This is a\nmultiline summary", rs[0].Summary())
}

/*
2020-01-15 (7h30m!)
This is a
multiline summary
	8:00 - 12:15
	2h15m
	14:38 -
	-1h51m
	<23:23 - 4:03
	22:00 - 0:01>
*/
