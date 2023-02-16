package cli

import (
	"github.com/jotaen/klog/klog"
	"github.com/jotaen/klog/klog/app/cli/lib"
	"github.com/jotaen/klog/klog/service"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestStartWithAutoTime(t *testing.T) {
	state, err := NewTestingContext()._SetRecords(`
1920-02-02
	9:00-12:00
`)._SetNow(1920, 2, 2, 15, 24)._Run((&Start{}).Run)
	require.Nil(t, err)
	assert.Equal(t, `
1920-02-02
	9:00-12:00
	15:24-?
`, state.writtenFileContents)
}

func TestStartWithExplicitDateAndAutoTimeYesterday(t *testing.T) {
	state, err := NewTestingContext()._SetRecords(`
1920-02-02
	9:00-12:00
`)._SetNow(1920, 2, 3, 23, 35)._Run((&Start{
		AtDateAndTimeArgs: lib.AtDateAndTimeArgs{
			AtDateArgs: lib.AtDateArgs{Date: klog.Ɀ_Date_(1920, 2, 2)},
		},
	}).Run)
	require.Nil(t, err)
	assert.Equal(t, `
1920-02-02
	9:00-12:00
	23:35>-?
`, state.writtenFileContents)
}

func TestStartWithExplicitTime(t *testing.T) {
	state, err := NewTestingContext()._SetRecords(`
1920-02-02
	9:00-12:00
`)._SetNow(1920, 2, 2, 23, 0)._Run((&Start{
		AtDateAndTimeArgs: lib.AtDateAndTimeArgs{
			Time: klog.Ɀ_Time_(15, 24),
		},
	}).Run)
	require.Nil(t, err)
	assert.Equal(t, `
1920-02-02
	9:00-12:00
	15:24-?
`, state.writtenFileContents)
}

func TestStartWithExplicitDateAndTime(t *testing.T) {
	state, err := NewTestingContext()._SetRecords(`
1920-02-02
	9:00-12:00
`)._SetNow(1920, 9, 28, 12, 16)._Run((&Start{
		AtDateAndTimeArgs: lib.AtDateAndTimeArgs{
			AtDateArgs: lib.AtDateArgs{Date: klog.Ɀ_Date_(1920, 2, 2)},
			Time:       klog.Ɀ_Time_(15, 24),
		},
	}).Run)
	require.Nil(t, err)
	assert.Equal(t, `
1920-02-02
	9:00-12:00
	15:24-?
`, state.writtenFileContents)
}

func TestStartFailsIfDateIsInPastAndNoTimeIsGiven(t *testing.T) {
	state, err := NewTestingContext()._SetRecords(`
1920-02-02
	9:00-???
`)._SetNow(1920, 9, 28, 12, 15)._Run((&Start{
		AtDateAndTimeArgs: lib.AtDateAndTimeArgs{
			AtDateArgs: lib.AtDateArgs{Date: klog.Ɀ_Date_(1920, 2, 2)},
		},
	}).Run)
	require.Error(t, err)
	assert.Equal(t, "Please specify a time value for dates in the past", err.Details())
	assert.Equal(t, state.writtenFileContents, "")
}

func TestStartFailsIfAlreadyStarted(t *testing.T) {
	state, err := NewTestingContext()._SetRecords(`
1920-02-02
	9:00-???
`)._Run((&Start{
		AtDateAndTimeArgs: lib.AtDateAndTimeArgs{
			AtDateArgs: lib.AtDateArgs{Date: klog.Ɀ_Date_(1920, 2, 2)},
			Time:       klog.Ɀ_Time_(12, 35),
		},
	}).Run)
	require.Error(t, err)
	assert.Equal(t, "There is already an open range in this record", err.Details())
	assert.Equal(t, state.writtenFileContents, "")
}

func TestStartWithSummary(t *testing.T) {
	state, err := NewTestingContext()._SetRecords(`
1920-02-02
	9:00-12:00
`)._SetNow(1920, 2, 2, 15, 24)._Run((&Start{
		AtDateAndTimeArgs: lib.AtDateAndTimeArgs{
			AtDateArgs: lib.AtDateArgs{Date: klog.Ɀ_Date_(1920, 2, 2)},
		},
		Summary: klog.Ɀ_EntrySummary_("Started something"),
	}).Run)
	require.Nil(t, err)
	assert.Equal(t, `
1920-02-02
	9:00-12:00
	15:24-? Started something
`, state.writtenFileContents)
}

func TestStartAtUnknownDateCreatesNewRecord(t *testing.T) {
	state, err := NewTestingContext()._SetRecords(`1623-12-13
	09:23 - ???
`)._SetNow(1623, 12, 11, 12, 49)._Run((&Start{}).Run)
	require.Nil(t, err)
	assert.Equal(t, `1623-12-11
	12:49 - ???

1623-12-13
	09:23 - ???
`, state.writtenFileContents)
}

func TestStartNewRecordWithShouldTotal(t *testing.T) {
	state, err := NewTestingContext()._SetRecords(`1623-12-13
	09:23 - ???
`)._SetNow(1623, 12, 11, 12, 49)._SetFileConfig(`
default_should_total = 8h!
`)._Run((&Start{}).Run)
	require.Nil(t, err)
	assert.Equal(t, `1623-12-11 (8h!)
	12:49 - ???

1623-12-13
	09:23 - ???
`, state.writtenFileContents)
}

func TestStartWithStyle(t *testing.T) {
	// For empty file and no preferences, use recommended default.
	{
		state, err := NewTestingContext()._SetRecords(`
1920/02/02
`)._SetNow(1920, 2, 2, 9, 44)._Run((&Start{}).Run)
		require.Nil(t, err)
		assert.Equal(t, `
1920/02/02
    9:44 - ?
`, state.writtenFileContents)
	}

	// Without any preference, detect from file.
	{
		state, err := NewTestingContext()._SetRecords(`
1920/02/02
  9:00am-1:00pm
  3h
`)._SetNow(1920, 2, 3, 8, 12)._Run((&Start{}).Run)
		require.Nil(t, err)
		assert.Equal(t, `
1920/02/02
  9:00am-1:00pm
  3h

1920/02/03
  8:12am-?
`, state.writtenFileContents)
	}

	// Use preference from config file, if given.
	{
		state, err := NewTestingContext()._SetRecords(`
1920/02/02
  9:00am-1:00pm
  3h
`)._SetNow(1920, 2, 3, 8, 12)._SetFileConfig(`
date_format = YYYY-MM-DD
time_format = 24h
`)._Run((&Start{}).Run)
		require.Nil(t, err)
		assert.Equal(t, `
1920/02/02
  9:00am-1:00pm
  3h

1920-02-03
  8:12-?
`, state.writtenFileContents)
	}

	// If explicit flag was provided, that takes ultimate precedence.
	{
		state, err := NewTestingContext()._SetRecords(`
1920/02/02
  9:00am-1:00pm
`)._SetFileConfig(`
time_format = 12h
`)._SetNow(1920, 2, 2, 8, 12)._Run((&Start{
			AtDateAndTimeArgs: lib.AtDateAndTimeArgs{
				Time: klog.Ɀ_Time_(9, 44),
			},
		}).Run)
		require.Nil(t, err)
		assert.Equal(t, `
1920/02/02
  9:00am-1:00pm
  9:44-?
`, state.writtenFileContents)
	}
}

func TestStartWithRounding(t *testing.T) {
	// With --round flag
	{
		r5, _ := service.NewRounding(5)
		state, err := NewTestingContext()._SetRecords(`
2005-05-05
`)._SetNow(2005, 5, 5, 8, 12)._Run((&Start{
			AtDateAndTimeArgs: lib.AtDateAndTimeArgs{Round: r5},
		}).Run)
		require.Nil(t, err)
		assert.Equal(t, `
2005-05-05
    8:10 - ?
`, state.writtenFileContents)
	}

	// With file config
	{
		state, err := NewTestingContext()._SetRecords(`
2005-05-05
`)._SetNow(2005, 5, 5, 8, 12)._SetFileConfig(`
default_rounding = 15m
`)._Run((&Start{}).Run)
		require.Nil(t, err)
		assert.Equal(t, `
2005-05-05
    8:15 - ?
`, state.writtenFileContents)
	}

	// Flag trumps file config
	{
		r5, _ := service.NewRounding(5)
		state, err := NewTestingContext()._SetRecords(`
2005-05-05
`)._SetNow(2005, 5, 5, 8, 12)._SetFileConfig(`
default_rounding = 60m
`)._Run((&Start{
			AtDateAndTimeArgs: lib.AtDateAndTimeArgs{Round: r5},
		}).Run)
		require.Nil(t, err)
		assert.Equal(t, `
2005-05-05
    8:10 - ?
`, state.writtenFileContents)
	}
}
