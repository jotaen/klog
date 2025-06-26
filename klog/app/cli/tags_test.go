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
	ctx := NewTestingContext()._SetRecords(`
1995-03-17
#sports
	3h #badminton
	1h #running=home-trail
	1h #running=river-route

1995-03-28
Was #sick, need to compensate later
	-30m #running

1995-04-02
	9h something untagged
	45m #badminton

1995-04-19
#sports #running (Don’t count that twice!)
	14:00 - 17:00 #sports #running
	
`)

	t.Run("Without argument", func(t *testing.T) {
		state, err := ctx._Run((&Tags{}).Run)
		require.Nil(t, err)
		assert.Equal(t, `
#badminton 3h45m
#running   4h30m
#sick      -30m 
#sports    8h   
`, state.printBuffer)
	})

	t.Run("With count", func(t *testing.T) {
		state, err := ctx._Run((&Tags{
			Count: true,
		}).Run)
		require.Nil(t, err)
		assert.Equal(t, `
#badminton 3h45m  (2)
#running   4h30m  (4)
#sick      -30m   (1)
#sports    8h     (4)
`, state.printBuffer)
	})

	t.Run("With values", func(t *testing.T) {
		state, err := ctx._Run((&Tags{
			Values: true,
		}).Run)
		require.Nil(t, err)
		assert.Equal(t, `
#badminton   3h45m   
#running     4h30m   
 home-trail        1h
 river-route       1h
#sick        -30m    
#sports      8h      
`, state.printBuffer)
	})

	t.Run("With values and count", func(t *testing.T) {
		state, err := ctx._Run((&Tags{
			Values: true,
			Count:  true,
		}).Run)
		require.Nil(t, err)
		assert.Equal(t, `
#badminton   3h45m     (2)
#running     4h30m     (4)
 home-trail        1h  (1)
 river-route       1h  (1)
#sick        -30m      (1)
#sports      8h        (4)
`, state.printBuffer)
	})

	t.Run("With untagged", func(t *testing.T) {
		state, err := ctx._Run((&Tags{
			WithUntagged: true,
		}).Run)
		require.Nil(t, err)
		assert.Equal(t, `
#badminton 3h45m
#running   4h30m
#sick      -30m 
#sports    8h   
(untagged) 9h   
`, state.printBuffer)
	})

	t.Run("With untagged and count", func(t *testing.T) {
		state, err := ctx._Run((&Tags{
			WithUntagged: true,
			Count:        true,
		}).Run)
		require.Nil(t, err)
		assert.Equal(t, `
#badminton 3h45m  (2)
#running   4h30m  (4)
#sick      -30m   (1)
#sports    8h     (4)
(untagged) 9h     (1)
`, state.printBuffer)
	})

	t.Run("With values and untagged", func(t *testing.T) {
		state, err := ctx._Run((&Tags{
			Values:       true,
			WithUntagged: true,
		}).Run)
		require.Nil(t, err)
		assert.Equal(t, `
#badminton   3h45m   
#running     4h30m   
 home-trail        1h
 river-route       1h
#sick        -30m    
#sports      8h      
(untagged)   9h      
`, state.printBuffer)
	})

	t.Run("With values and untagged and count", func(t *testing.T) {
		state, err := ctx._Run((&Tags{
			Values:       true,
			WithUntagged: true,
			Count:        true,
		}).Run)
		require.Nil(t, err)
		assert.Equal(t, `
#badminton   3h45m     (2)
#running     4h30m     (4)
 home-trail        1h  (1)
 river-route       1h  (1)
#sick        -30m      (1)
#sports      8h        (4)
(untagged)   9h        (1)
`, state.printBuffer)
	})
}

func TestPrintUntaggedIfNoTags(t *testing.T) {
	t.Run("No tags present", func(t *testing.T) {
		state, err := NewTestingContext()._SetRecords(`
1995-03-17
	1h
`)._Run((&Tags{
			WithUntagged: true,
		}).Run)
		require.Nil(t, err)
		assert.Equal(t, `
(untagged) 1h
`, state.printBuffer)
	})

	t.Run("Empty file", func(t *testing.T) {
		state, err := NewTestingContext()._SetRecords(`
`)._Run((&Tags{
			WithUntagged: true,
		}).Run)
		require.Nil(t, err)
		assert.Equal(t, `
(untagged) 0m
`, state.printBuffer)
	})

	t.Run("Empty file (with count)", func(t *testing.T) {
		state, err := NewTestingContext()._SetRecords(`
`)._Run((&Tags{
			WithUntagged: true,
			Count:        true,
		}).Run)
		require.Nil(t, err)
		assert.Equal(t, `
(untagged) 0m  (0)
`, state.printBuffer)
	})
}

func TestPrintTagsWithUnicodeCharacters(t *testing.T) {
	state, err := NewTestingContext()._SetRecords(`
1995-03-17
	1h #ascii
	2h #üñïčöδę
`)._Run((&Tags{}).Run)
	require.Nil(t, err)
	assert.Equal(t, `
#ascii   1h
#üñïčöδę 2h
`, state.printBuffer)
}

func TestPrintTagsOverviewWithUntaggedEmptyStates(t *testing.T) {
	ctx := NewTestingContext()._SetRecords(`
1995-03-17
	3h #ticket
`)
	t.Run("Include 0 line", func(t *testing.T) {
		state, err := ctx._Run((&Tags{
			WithUntagged: true,
		}).Run)
		require.Nil(t, err)
		assert.Equal(t, `
#ticket    3h
(untagged) 0m
`, state.printBuffer)
	})

	t.Run("Include 0 count", func(t *testing.T) {
		state, err := ctx._Run((&Tags{
			WithUntagged: true,
			Count:        true,
		}).Run)
		require.Nil(t, err)
		assert.Equal(t, `
#ticket    3h  (1)
(untagged) 0m  (0)
`, state.printBuffer)
	})
}

func TestPrintTagsOverviewWithValueGrouping(t *testing.T) {
	state, err := NewTestingContext()._SetRecords(`
1995-03-17
	3h #ticket=481
	1h #ticket=105
	1h
`)._Run((&Tags{Values: true}).Run)
	require.Nil(t, err)
	assert.Equal(t, `
#ticket 4h   
 105       1h
 481       3h
`, state.printBuffer)
}

func TestPrintTagsOverviewWithValueGroupingAndCount(t *testing.T) {
	state, err := NewTestingContext()._SetRecords(`
1995-03-17
	3h #ticket=481
	1h #ticket=105
	1h
`)._Run((&Tags{
		Values: true,
		Count:  true,
	}).Run)
	require.Nil(t, err)
	assert.Equal(t, `
#ticket 4h     (2)
 105       1h  (1)
 481       3h  (1)
`, state.printBuffer)
}

func TestPrintTagsOverviewWithoutValueGrouping(t *testing.T) {
	state, err := NewTestingContext()._SetRecords(`
1995-03-17
	3h #ticket=481
	1h #ticket=105
	1h
`)._Run((&Tags{}).Run)
	require.Nil(t, err)
	assert.Equal(t, `
#ticket 4h
`, state.printBuffer)
}
