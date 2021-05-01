package cli

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"klog"
	"klog/app/cli/lib"
	"testing"
)

func TestTrackEntry(t *testing.T) {
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
