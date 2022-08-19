package cli

import (
	"github.com/jotaen/klog/klog"
	"github.com/jotaen/klog/klog/app"
	"github.com/jotaen/klog/klog/app/cli/lib"
	"github.com/jotaen/klog/klog/service"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestStop(t *testing.T) {
	state, err := NewTestingContext()._SetRecords(`
1920-02-02
	11:22-?
`)._SetNow(1920, 2, 2, 15, 24)._Run((&Stop{
		AtDateAndTimeArgs: lib.AtDateAndTimeArgs{
			AtDateArgs: lib.AtDateArgs{Date: klog.Ɀ_Date_(1920, 2, 2)},
		},
	}).Run)
	require.Nil(t, err)
	assert.Equal(t, `
1920-02-02
	11:22-15:24
`, state.writtenFileContents)
}

func TestStopFallsBackWithShiftedTimeToYesterdayWithAutoTime(t *testing.T) {
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

func TestDoesNotFallBackToYesterdayWhenTimeIsExplicit(t *testing.T) {
	_, err := NewTestingContext()._SetRecords(`
1920-02-02
	22:22-?
`)._SetNow(1920, 2, 3, 4, 16)._Run((&Stop{
		AtDateAndTimeArgs: lib.AtDateAndTimeArgs{Time: klog.Ɀ_Time_(23, 30)},
	}).Run)
	require.Error(t, err)
}

func TestStopsYesterdaysRecordWithShiftedAutoTime(t *testing.T) {
	state, err := NewTestingContext()._SetRecords(`
1920-02-02
	10:22pm-?
`)._SetNow(1920, 2, 3, 2, 49)._Run((&Stop{
		AtDateAndTimeArgs: lib.AtDateAndTimeArgs{
			AtDateArgs: lib.AtDateArgs{Date: klog.Ɀ_Date_(1920, 2, 2)},
		},
	}).Run)
	require.Nil(t, err)
	assert.Equal(t, `
1920-02-02
	10:22pm-2:49am>
`, state.writtenFileContents)
}

func TestStopWithExtendingSummary(t *testing.T) {
	state, err := NewTestingContext()._SetRecords(`
1920-02-02
	11:22-? Started something...
`)._SetNow(1920, 2, 2, 15, 24)._Run((&Stop{
		AtDateAndTimeArgs: lib.AtDateAndTimeArgs{
			AtDateArgs: lib.AtDateArgs{Date: klog.Ɀ_Date_(1920, 2, 2)},
		},
		Summary: klog.Ɀ_EntrySummary_("Done!"),
	}).Run)
	require.Nil(t, err)
	assert.Equal(t, `
1920-02-02
	11:22-15:24 Started something... Done!
`, state.writtenFileContents)
}

func TestStopFailsIfNoTimeSpecifiedForPastDates(t *testing.T) {
	state, err := NewTestingContext()._SetRecords(`
1623-12-12
	15:00-?
`)._Run((&Stop{
		AtDateAndTimeArgs: lib.AtDateAndTimeArgs{
			AtDateArgs: lib.AtDateArgs{Date: klog.Ɀ_Date_(1624, 02, 1)},
		},
	}).Run)
	require.Error(t, err)
	assert.Equal(t, "Missing time parameter", err.Error())
	assert.Equal(t, "", state.writtenFileContents)
}

func TestStopFailsIfNoRecentRecord(t *testing.T) {
	state, err := NewTestingContext()._SetRecords(`
1623-12-12
	15:00-?
`)._Run((&Stop{
		AtDateAndTimeArgs: lib.AtDateAndTimeArgs{
			AtDateArgs: lib.AtDateArgs{Date: klog.Ɀ_Date_(1624, 02, 1)},
			Time:       klog.Ɀ_Time_(16, 00),
		},
	}).Run)
	require.Error(t, err)
	assert.Equal(t, "No such record", err.Error())
	assert.Equal(t, "", state.writtenFileContents)
}

func TestStopFailsIfNoOpenRange(t *testing.T) {
	state, err := NewTestingContext()._SetRecords(`
1623-12-12
	15:00-16:00
`)._Run((&Stop{
		AtDateAndTimeArgs: lib.AtDateAndTimeArgs{
			AtDateArgs: lib.AtDateArgs{Date: klog.Ɀ_Date_(1623, 12, 12)},
			Time:       klog.Ɀ_Time_(16, 00),
		},
	}).Run)
	require.Error(t, err)
	assert.Equal(t, "Manipulation failed", err.(app.Error).Error())
	assert.Equal(t, "No open time range", err.(app.Error).Details())
	assert.Equal(t, "", state.writtenFileContents)
}

func TestStopWithRounding(t *testing.T) {
	r15, _ := service.NewRounding(15)
	state, err := NewTestingContext()._SetRecords(`
2005-05-05
    8:10 - ?
`)._SetNow(2005, 5, 5, 11, 24)._Run((&Stop{
		AtDateAndTimeArgs: lib.AtDateAndTimeArgs{Round: r15},
	}).Run)
	require.Nil(t, err)
	assert.Equal(t, `
2005-05-05
    8:10 - 11:30
`, state.writtenFileContents)
}
