package cli

import (
	"testing"

	"github.com/jotaen/klog/klog"
	"github.com/jotaen/klog/klog/app/cli/args"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestTrackEntryInEmptyFile(t *testing.T) {
	state, err := NewTestingContext()._SetRecords("")._Run((&Track{
		Entry:      klog.Ɀ_EntrySummary_("2h"),
		AtDateArgs: args.AtDateArgs{Date: klog.Ɀ_Date_(1855, 4, 25)},
	}).Run)
	require.Nil(t, err)
	assert.Equal(t, "1855-04-25\n    2h\n", state.writtenFileContents)
}

func TestTrackEntryInExistingFile(t *testing.T) {
	state, err := NewTestingContext()._SetRecords(`
1855-04-25
	1h
`)._Run((&Track{
		Entry:      klog.Ɀ_EntrySummary_("2h"),
		AtDateArgs: args.AtDateArgs{Date: klog.Ɀ_Date_(1855, 4, 25)},
	}).Run)
	require.Nil(t, err)
	assert.Equal(t, `
1855-04-25
	1h
	2h
`, state.writtenFileContents)
}

func TestTrackEntryAtUnknownDateCreatesNewRecord(t *testing.T) {
	state, err := NewTestingContext()._SetRecords(`
1855-04-25
	1h
`)._Run((&Track{
		Entry:      klog.Ɀ_EntrySummary_("2h"),
		AtDateArgs: args.AtDateArgs{Date: klog.Ɀ_Date_(2000, 1, 1)},
	}).Run)
	require.Nil(t, err)
	assert.Equal(t, `
1855-04-25
	1h

2000-01-01
	2h
`, state.writtenFileContents)
}

func TestTrackNewRecordWithShouldTotal(t *testing.T) {
	state, err := NewTestingContext()._SetRecords(`
1855-04-25
	1h
`)._SetFileConfig(`
default_should_total = 7h30m!
`)._Run((&Track{
		Entry:      klog.Ɀ_EntrySummary_("2h"),
		AtDateArgs: args.AtDateArgs{Date: klog.Ɀ_Date_(2000, 1, 1)},
	}).Run)
	require.Nil(t, err)
	assert.Equal(t, `
1855-04-25
	1h

2000-01-01 (7h30m!)
	2h
`, state.writtenFileContents)
}

func TestTrackFailsIfEntryInvalid(t *testing.T) {
	state, err := NewTestingContext()._SetRecords(`
1855-04-25
	1h
`)._Run((&Track{
		Entry:      klog.Ɀ_EntrySummary_("Foo"),
		AtDateArgs: args.AtDateArgs{Date: klog.Ɀ_Date_(1855, 4, 25)},
	}).Run)
	require.Error(t, err)
	assert.Equal(t, "Manipulation failed", err.Error())
	assert.Equal(t, "This operation wouldn’t result in a valid record", err.Details())
	assert.Equal(t, "", state.writtenFileContents)
}

func TestTrackWithStyle(t *testing.T) {
	t.Run("For empty file and no preferences, use recommended default.", func(t *testing.T) {
		state, err := NewTestingContext()._SetRecords("").
			_SetNow(2000, 1, 1, 12, 00).
			_Run((&Track{
				Entry: klog.Ɀ_EntrySummary_("2h"),
			}).Run)
		require.Nil(t, err)
		assert.Equal(t, "2000-01-01\n    2h\n", state.writtenFileContents)
	})

	t.Run("Without any preference, detect from file.", func(t *testing.T) {
		state, err := NewTestingContext()._SetRecords(`
1855/04/25
	1h
`)._SetNow(2000, 1, 1, 12, 00)._Run((&Track{
			Entry: klog.Ɀ_EntrySummary_("2h"),
		}).Run)
		require.Nil(t, err)
		assert.Equal(t, `
1855/04/25
	1h

2000/01/01
	2h
`, state.writtenFileContents)
	})

	t.Run("Use preference from config file, if given.", func(t *testing.T) {
		state, err := NewTestingContext()._SetRecords(`
1855/04/25
	1h
`)._SetFileConfig(`
date_format = YYYY-MM-DD
`)._SetNow(2000, 1, 1, 12, 00)._Run((&Track{
			Entry: klog.Ɀ_EntrySummary_("2h"),
		}).Run)
		require.Nil(t, err)
		assert.Equal(t, `
1855/04/25
	1h

2000-01-01
	2h
`, state.writtenFileContents)
	})

	t.Run("If explicit flag was provided, that takes ultimate precedence.", func(t *testing.T) {
		state, err := NewTestingContext()._SetRecords(`
1855/04/25
	1h
`)._SetFileConfig(`
date_format = YYYY/MM/DD
`)._Run((&Track{
			Entry:      klog.Ɀ_EntrySummary_("2h"),
			AtDateArgs: args.AtDateArgs{Date: klog.Ɀ_Date_(2000, 1, 1)},
		}).Run)
		require.Nil(t, err)
		assert.Equal(t, `
1855/04/25
	1h

2000-01-01
	2h
`, state.writtenFileContents)
	})
}
