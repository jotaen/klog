package commands

import (
	"github.com/jotaen/klog/src"
	"github.com/jotaen/klog/src/app"
	"github.com/jotaen/klog/src/app/lib"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestStop(t *testing.T) {
	state, err := NewTestingContext()._SetRecords(`
1920-02-02
	11:22-?
`)._SetNow(1920, 2, 2, 15, 24)._Run((&Stop{
		AtDateArgs: lib.AtDateArgs{Date: klog.Ɀ_Date_(1920, 2, 2)},
	}).Run)
	require.Nil(t, err)
	assert.Equal(t, `
1920-02-02
	11:22-15:24
`, state.writtenFileContents)
}

func TestStopFallsBackToYesterday(t *testing.T) {
	state, err := NewTestingContext()._SetRecords(`
1920-02-02
	22:22-?
`)._SetNow(1920, 2, 3, 4, 16)._Run((&Stop{}).Run)
	require.Nil(t, err)
	assert.Equal(t, `
1920-02-02
	22:22-4:16>
`, state.writtenFileContents)
}

func TestStopWithExtendingSummary(t *testing.T) {
	state, err := NewTestingContext()._SetRecords(`
1920-02-02
	11:22-? Started something...
`)._SetNow(1920, 2, 2, 15, 24)._Run((&Stop{
		AtDateArgs: lib.AtDateArgs{Date: klog.Ɀ_Date_(1920, 2, 2)},
		Summary:    "Done!",
	}).Run)
	require.Nil(t, err)
	assert.Equal(t, `
1920-02-02
	11:22-15:24 Started something... Done!
`, state.writtenFileContents)
}

func TestStopFailsIfNoRecentRecord(t *testing.T) {
	state, err := NewTestingContext()._SetRecords(`
1623-12-12
	15:00-?
`)._Run((&Stop{
		AtDateArgs: lib.AtDateArgs{Date: klog.Ɀ_Date_(1624, 02, 1)},
	}).Run)
	require.Error(t, err)
	assert.Equal(t, err.Error(), "No such record")
	assert.Equal(t, state.writtenFileContents, "")
}

func TestStopFailsIfNoOpenRange(t *testing.T) {
	state, err := NewTestingContext()._SetRecords(`
1623-12-12
	15:00-16:00
`)._Run((&Stop{
		AtDateArgs: lib.AtDateArgs{Date: klog.Ɀ_Date_(1623, 12, 12)},
	}).Run)
	require.Error(t, err)
	assert.Equal(t, "Manipulation failed", err.(app.Error).Error())
	assert.Equal(t, "No open time range", err.(app.Error).Details())
	assert.Equal(t, state.writtenFileContents, "")
}
