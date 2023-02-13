package cli

import (
	"github.com/jotaen/klog/klog"
	"github.com/jotaen/klog/klog/app/cli/lib"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestTrackEntryInEmptyFile(t *testing.T) {
	state, err := NewTestingContext()._SetRecords("")._Run((&Track{
		Entry:      klog.Ɀ_EntrySummary_("2h"),
		AtDateArgs: lib.AtDateArgs{Date: klog.Ɀ_Date_(1855, 4, 25)},
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
		AtDateArgs: lib.AtDateArgs{Date: klog.Ɀ_Date_(1855, 4, 25)},
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
		AtDateArgs: lib.AtDateArgs{Date: klog.Ɀ_Date_(2000, 1, 1)},
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
default_should_total: 7h30m!
`)._Run((&Track{
		Entry:      klog.Ɀ_EntrySummary_("2h"),
		AtDateArgs: lib.AtDateArgs{Date: klog.Ɀ_Date_(2000, 1, 1)},
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
		AtDateArgs: lib.AtDateArgs{Date: klog.Ɀ_Date_(1855, 4, 25)},
	}).Run)
	require.Error(t, err)
	assert.Equal(t, "Manipulation failed", err.Error())
	assert.Equal(t, "This operation wouldn’t result in a valid record", err.Details())
	assert.Equal(t, "", state.writtenFileContents)
}
