package cli

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"klog"
	"klog/app/cli/lib"
	"testing"
)

func TestStart(t *testing.T) {
	state, err := NewTestingContext()._SetRecords(`
1920-02-02
	9:00-12:00
`)._SetNow(1920, 2, 2, 15, 24)._Run((&Start{
		AtDateArgs: lib.AtDateArgs{Date: klog.Ɀ_Date_(1920, 2, 2)},
	}).Run)
	require.Nil(t, err)
	assert.Equal(t, `
1920-02-02
	9:00-12:00
	15:24 - ?
`, state.writtenFileContents)
}

func TestStartFailsIfAlreadyStarted(t *testing.T) {
	state, err := NewTestingContext()._SetRecords(`
1920-02-02
	9:00-???
`)._Run((&Start{
		AtDateArgs: lib.AtDateArgs{Date: klog.Ɀ_Date_(1920, 2, 2)},
	}).Run)
	require.Error(t, err)
	assert.Equal(t, state.writtenFileContents, "")
}

func TestStartWithSummary(t *testing.T) {
	state, err := NewTestingContext()._SetRecords(`
1920-02-02
	9:00-12:00
`)._SetNow(1920, 2, 2, 15, 24)._Run((&Start{
		AtDateArgs: lib.AtDateArgs{Date: klog.Ɀ_Date_(1920, 2, 2)},
		Summary:    "Started something",
	}).Run)
	require.Nil(t, err)
	assert.Equal(t, `
1920-02-02
	9:00-12:00
	15:24 - ? Started something
`, state.writtenFileContents)
}

func TestStartAtUnknownDateCreatesNewRecord(t *testing.T) {
	state, err := NewTestingContext()._SetRecords(`1623-12-13
	12:23-???
`)._SetNow(1623, 12, 11, 12, 49)._Run((&Start{}).Run)
	require.Nil(t, err)
	assert.Equal(t, `1623-12-11
	12:49 - ?

1623-12-13
	12:23-???
`, state.writtenFileContents)
}
