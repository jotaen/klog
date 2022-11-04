package reconciling

import (
	"github.com/jotaen/klog/klog"
	"github.com/jotaen/klog/klog/parser"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestReconcilerAddsNewPauseEntry(t *testing.T) {
	original := `
2010-04-27
    3:00pm - ?
`
	rs, _ := parser.NewSerialParser().Parse(original)
	reconciler := NewReconcilerAtRecord(rs, klog.Ɀ_Date_(2010, 4, 27))
	require.NotNil(t, reconciler)
	result, err := reconciler.PauseOpenRange(klog.NewDuration(0, -12), nil)
	require.Nil(t, err)
	assert.Equal(t, `
2010-04-27
    3:00pm - ?
    -12m
`, result.AllSerialised)
}

func TestReconcilerFailsIfPauseIsPositiveValue(t *testing.T) {
	original := `
2010-04-27
    3:00 - 4:00
`
	rs, _ := parser.NewSerialParser().Parse(original)
	reconciler := NewReconcilerAtRecord(rs, klog.Ɀ_Date_(2010, 4, 27))
	require.NotNil(t, reconciler)
	result, err := reconciler.PauseOpenRange(klog.NewDuration(0, 12), nil)
	require.Error(t, err)
	assert.Nil(t, result)
}

func TestReconcilerFailsIfThereIsNoOpenRange(t *testing.T) {
	original := `
2010-04-27
    3:00 - 4:00
`
	rs, _ := parser.NewSerialParser().Parse(original)
	reconciler := NewReconcilerAtRecord(rs, klog.Ɀ_Date_(2010, 4, 27))
	require.NotNil(t, reconciler)
	result, err := reconciler.PauseOpenRange(klog.NewDuration(0, -12), nil)
	require.Error(t, err)
	assert.Nil(t, result)
}

func TestReconcilerAddsToExistingPauseEntry(t *testing.T) {
	original := `
2010-04-27
Foo
    3:00 - ? I desperately need
        a break!
    -30m
`
	rs, _ := parser.NewSerialParser().Parse(original)
	reconciler := NewReconcilerAtRecord(rs, klog.Ɀ_Date_(2010, 4, 27))
	require.NotNil(t, reconciler)
	result, err := reconciler.PauseOpenRange(klog.NewDuration(0, -3), nil)
	require.Nil(t, err)
	assert.Equal(t, `
2010-04-27
Foo
    3:00 - ? I desperately need
        a break!
    -33m
`, result.AllSerialised)
}

func TestReconcilerOnlyAddsToExistingPauseEntryIfSummaryMatches(t *testing.T) {
	original := `
2010-04-27
Foo
    3:00 - ?
    -30m This is a totally unrelated entry,
        that should not be modified!
`
	rs, _ := parser.NewSerialParser().Parse(original)
	reconciler := NewReconcilerAtRecord(rs, klog.Ɀ_Date_(2010, 4, 27))
	require.NotNil(t, reconciler)
	result, err := reconciler.PauseOpenRange(klog.NewDuration(-1, -30), klog.Ɀ_EntrySummary_("Lunch break"))
	require.Nil(t, err)
	assert.Equal(t, `
2010-04-27
Foo
    3:00 - ?
    -1h30m Lunch break
    -30m This is a totally unrelated entry,
        that should not be modified!
`, result.AllSerialised)
}

func TestReconcilerDoesNotExtendNonNegativeDurations(t *testing.T) {
	original := `
2010-04-27
    3:00 - ?
    30m
`
	rs, _ := parser.NewSerialParser().Parse(original)
	reconciler := NewReconcilerAtRecord(rs, klog.Ɀ_Date_(2010, 4, 27))
	require.NotNil(t, reconciler)
	result, err := reconciler.PauseOpenRange(klog.NewDuration(0, -10), nil)
	require.Nil(t, err)
	assert.Equal(t, `
2010-04-27
    3:00 - ?
    -10m
    30m
`, result.AllSerialised)
}
