package cli

import (
	"github.com/jotaen/klog/klog"
	"github.com/jotaen/klog/klog/app/cli/util"
	"github.com/jotaen/klog/klog/service"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestSwitch(t *testing.T) {
	state, err := NewTestingContext()._SetRecords(`
1920-02-02
	11:22-?
`)._SetNow(1920, 2, 2, 15, 24)._Run((&Switch{
		AtDateAndTimeArgs: util.AtDateAndTimeArgs{
			AtDateArgs: util.AtDateArgs{Date: klog.Ɀ_Date_(1920, 2, 2)},
		},
	}).Run)
	require.Nil(t, err)
	assert.Equal(t, `
1920-02-02
	11:22-15:24
	15:24-?
`, state.writtenFileContents)
}

func TestSwitchWithSummaries(t *testing.T) {
	state, err := NewTestingContext()._SetRecords(`
1920-02-02
	11:22-? Currently ongoing...
		...task

1920-02-03
Next day
`)._SetNow(1920, 2, 2, 15, 24)._Run((&Switch{
		AtDateAndTimeArgs: util.AtDateAndTimeArgs{
			AtDateArgs: util.AtDateArgs{Date: klog.Ɀ_Date_(1920, 2, 2)},
		},
		SummaryArgs: util.SummaryArgs{
			SummaryText: klog.Ɀ_EntrySummary_("Start", "over"),
		},
	}).Run)
	require.Nil(t, err)
	assert.Equal(t, `
1920-02-02
	11:22-15:24 Currently ongoing...
		...task
	15:24-? Start
		over

1920-02-03
Next day
`, state.writtenFileContents)
}

func TestSwitchWithResume(t *testing.T) {
	state, err := NewTestingContext()._SetRecords(`
1920-02-03
	8:00 - 9:00 First
	9:00 - ? Second
`)._SetNow(1920, 2, 3, 9, 31)._Run((&Switch{
		SummaryArgs: util.SummaryArgs{
			Resume: true,
		},
	}).Run)
	require.Nil(t, err)
	assert.Equal(t, `
1920-02-03
	8:00 - 9:00 First
	9:00 - 9:31 Second
	9:31 - ? Second
`, state.writtenFileContents)
}

func TestSwitchWithResumeNth(t *testing.T) {
	state, err := NewTestingContext()._SetRecords(`
1920-02-03
	8:00 - 9:00 First
	9:00 - ? Second
`)._SetNow(1920, 2, 3, 9, 31)._Run((&Switch{
		SummaryArgs: util.SummaryArgs{
			ResumeNth: 1,
		},
	}).Run)
	require.Nil(t, err)
	assert.Equal(t, `
1920-02-03
	8:00 - 9:00 First
	9:00 - 9:31 Second
	9:31 - ? First
`, state.writtenFileContents)
}

func TestSwitchCannotResumeAndSummary(t *testing.T) {
	state, err := NewTestingContext()._SetRecords(`
1920-02-03
	8:00 - 9:00 First
	9:00 - ? Second
`)._SetNow(1920, 2, 3, 9, 31)._Run((&Switch{
		SummaryArgs: util.SummaryArgs{
			Resume:      true,
			SummaryText: klog.Ɀ_EntrySummary_("Foo"),
		},
	}).Run)
	require.Error(t, err)
	assert.Equal(t, "Manipulation failed", err.Error())
	assert.Equal(t, "", state.writtenFileContents)
}

func TestSwitchWithStyle(t *testing.T) {
	t.Run("Without any preference, detect from file.", func(t *testing.T) {
		state, err := NewTestingContext()._SetRecords(`
1920-02-02
	10:22am-???
`)._SetNow(1920, 2, 2, 14, 49)._Run((&Switch{}).Run)
		require.Nil(t, err)
		assert.Equal(t, `
1920-02-02
	10:22am-2:49pm
	2:49pm-???
`, state.writtenFileContents)
	})

	t.Run("Use preference from config file, if given.", func(t *testing.T) {
		state, err := NewTestingContext()._SetRecords(`
1920-02-02
	10:22am-?
`)._SetFileConfig(`
time_convention = 24h
`)._SetNow(1920, 2, 2, 14, 49)._Run((&Switch{}).Run)
		require.Nil(t, err)
		assert.Equal(t, `
1920-02-02
	10:22am-14:49
	14:49-?
`, state.writtenFileContents)
	})

	t.Run("If explicit flag was provided, that takes ultimate precedence.", func(t *testing.T) {
		state, err := NewTestingContext()._SetRecords(`
1920-02-02
	10:22am-?
`)._SetFileConfig(`
time_convention = 12h
`)._SetNow(1920, 2, 2, 14, 49)._Run((&Switch{
			AtDateAndTimeArgs: util.AtDateAndTimeArgs{Time: klog.Ɀ_Time_(14, 49)},
		}).Run)
		require.Nil(t, err)
		assert.Equal(t, `
1920-02-02
	10:22am-14:49
	14:49-?
`, state.writtenFileContents)
	})
}

func TestSwitchFailsIfNoRecentRecord(t *testing.T) {
	state, err := NewTestingContext()._SetRecords(`
1623-12-12
	15:00-?
`)._Run((&Switch{
		AtDateAndTimeArgs: util.AtDateAndTimeArgs{
			AtDateArgs: util.AtDateArgs{Date: klog.Ɀ_Date_(1624, 02, 1)},
			Time:       klog.Ɀ_Time_(16, 00),
		},
	}).Run)
	require.Error(t, err)
	assert.Equal(t, "No such record", err.Error())
	assert.Equal(t, "", state.writtenFileContents)
}

func TestSwitchFailsIfNoOpenRange(t *testing.T) {
	state, err := NewTestingContext()._SetRecords(`
1623-12-12
	15:00-16:00
`)._Run((&Switch{
		AtDateAndTimeArgs: util.AtDateAndTimeArgs{
			AtDateArgs: util.AtDateArgs{Date: klog.Ɀ_Date_(1623, 12, 12)},
			Time:       klog.Ɀ_Time_(16, 00),
		},
	}).Run)
	require.Error(t, err)
	assert.Equal(t, "Manipulation failed", err.Error())
	assert.Equal(t, "No open time range", err.Details())
	assert.Equal(t, "", state.writtenFileContents)
}

func TestSwitchWithRounding(t *testing.T) {
	t.Run("With --round flag", func(t *testing.T) {
		r15, _ := service.NewRounding(15)
		state, err := NewTestingContext()._SetRecords(`
2005-05-05
    8:10 - ?
`)._SetNow(2005, 5, 5, 11, 24)._Run((&Switch{
			AtDateAndTimeArgs: util.AtDateAndTimeArgs{Round: r15},
		}).Run)
		require.Nil(t, err)
		assert.Equal(t, `
2005-05-05
    8:10 - 11:30
    11:30 - ?
`, state.writtenFileContents)
	})

	t.Run("With file config", func(t *testing.T) {
		state, err := NewTestingContext()._SetRecords(`
2005-05-05
    8:10 - ?
`)._SetNow(2005, 5, 5, 11, 24)._SetFileConfig(`
default_rounding = 30m
`)._Run((&Switch{}).Run)
		require.Nil(t, err)
		assert.Equal(t, `
2005-05-05
    8:10 - 11:30
    11:30 - ?
`, state.writtenFileContents)
	})

	t.Run("--round flag trumps file config", func(t *testing.T) {
		r5, _ := service.NewRounding(5)
		state, err := NewTestingContext()._SetRecords(`
2005-05-05
    8:10 - ?
`)._SetNow(2005, 5, 5, 11, 24)._SetFileConfig(`
default_rounding = 30m
`)._Run((&Switch{
			AtDateAndTimeArgs: util.AtDateAndTimeArgs{Round: r5},
		}).Run)
		require.Nil(t, err)
		assert.Equal(t, `
2005-05-05
    8:10 - 11:25
    11:25 - ?
`, state.writtenFileContents)
	})
}
