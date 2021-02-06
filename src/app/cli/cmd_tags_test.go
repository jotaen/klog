package cli

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestTagsOfEmptyInput(t *testing.T) {
	out, err := RunWithContext(``, (&Tags{}).Run)
	require.Nil(t, err)
	assert.Equal(t, "", out)
}

func TestPrintTagsOverview(t *testing.T) {
	/*
		Aspects tested:
		- Aggregate totals by tags
		- Sort output alphabetically
		- Print in tabular manner
	*/
	out, err := RunWithContext(`
1995-03-17
#sports
	3h #badminton
	1h #running
	1h #running

1995-03-28
Was #sick, need to compensate later
	-30m #running

1995-04-02
	99h something untagged
	45m #badminton

1995-04-19
#sports #running (Don’t count that twice!)
	14:00 - 17:00 #sports #running
	
`, (&Tags{}).Run)
	require.Nil(t, err)
	assert.Equal(t, `
#badminton 3h45m
#running   4h30m
#sick      -30m
#sports    8h
`, out)
}
