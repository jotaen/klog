package commands

import (
	"github.com/jotaen/klog/src"
	"github.com/jotaen/klog/src/app"
	"github.com/jotaen/klog/src/app/lib"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestTrackEntryInEmptyFile(t *testing.T) {
	state, err := NewTestingContext()._SetRecords("")._Run((&Track{
		Entry:      "2h",
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
		Entry:      "2h",
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
		Entry:      "2h",
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

func TestTrackFailsIfEntryInvalid(t *testing.T) {
	state, err := NewTestingContext()._SetRecords(`
1855-04-25
	1h
`)._Run((&Track{
		Entry:      "Foo",
		AtDateArgs: lib.AtDateArgs{Date: klog.Ɀ_Date_(1855, 4, 25)},
	}).Run)
	require.Error(t, err)
	assert.Equal(t, "Manipulation failed", err.Error())
	assert.Equal(t, "This operation wouldn’t result in a valid record", err.(app.Error).Details())
	assert.Equal(t, "", state.writtenFileContents)
}
