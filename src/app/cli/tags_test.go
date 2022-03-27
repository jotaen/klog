package cli

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestTagsOfEmptyInput(t *testing.T) {
	state, err := NewTestingContext()._SetRecords(``)._Run((&Tags{}).Run)
	require.Nil(t, err)
	assert.Equal(t, "", state.printBuffer)
}

func TestPrintTagsOverview(t *testing.T) {
	/*
		Aspects tested:
		- Aggregate totals by tags
		- Sort output alphabetically
		- Print in tabular manner
	*/
	state, err := NewTestingContext()._SetRecords(`
1995-03-17
#sports
	3h #badminton
	1h #running
	1h #running

1995-03-28
Was #sick, need to compensate later
	-30m #running

1995-04-02
	9h something untagged
	45m #badminton

1995-04-19
#sports #running (Don’t count that twice!)
	14:00 - 17:00 #sports #running
	
`)._Run((&Tags{}).Run)
	require.Nil(t, err)
	assert.Equal(t, `
#badminton 3h45m
#running   4h30m
#sick      -30m 
#sports    8h   
`, state.printBuffer)
}

func TestPrintTagsOverviewWithValueGrouping(t *testing.T) {
	state, err := NewTestingContext()._SetRecords(`
1995-03-17
	3h #ticket=481
	1h #ticket=105
	1h
`)._Run((&Tags{}).Run)
	require.Nil(t, err)
	assert.Equal(t, `
#ticket 4h
 105    1h
 481    3h
`, state.printBuffer)
}
