package reconciling

import (
	. "github.com/jotaen/klog/src"
	"github.com/jotaen/klog/src/parser"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestReconcilerClosesOpenRangeWithStyle(t *testing.T) {
	original := `
2010-04-27
    3:00pm - ??
`
	rs, _ := parser.Parse(original)
	reconciler := NewReconcilerAtRecord(rs, Ɀ_Date_(2010, 4, 27))
	require.NotNil(t, reconciler)
	result, err := reconciler.CloseOpenRange(Ɀ_Time_(15, 30), "")
	require.Nil(t, err)
	assert.Equal(t, `
2010-04-27
    3:00pm - 3:30pm
`, result.AllSerialised)
}

func TestReconcilerClosesOpenRangeWithNewSummary(t *testing.T) {
	original := `
2018-01-01
    15:00 - ?
`
	rs, _ := parser.Parse(original)
	reconciler := NewReconcilerAtRecord(rs, Ɀ_Date_(2018, 1, 1))
	require.NotNil(t, reconciler)
	result, err := reconciler.CloseOpenRange(Ɀ_Time_(15, 22), "Finished.")
	require.Nil(t, err)
	assert.Equal(t, `
2018-01-01
    15:00 - 15:22 Finished.
`, result.AllSerialised)
}

func TestReconcilerClosesOpenRangeWithNewMultilineSummary(t *testing.T) {
	original := `
2018-01-01
    15:00 - ?
`
	rs, _ := parser.Parse(original)
	reconciler := NewReconcilerAtRecord(rs, Ɀ_Date_(2018, 1, 1))
	require.NotNil(t, reconciler)
	result, err := reconciler.CloseOpenRange(Ɀ_Time_(15, 22), "\nFinished.")
	require.Nil(t, err)
	assert.Equal(t, `
2018-01-01
    15:00 - 15:22
        Finished.
`, result.AllSerialised)
}

func TestReconcilerClosesOpenRangeWithExtendingSummary(t *testing.T) {
	original := `
2018-01-01
    1h Multiline...
        ...entry summary
    15:00-??? Will this close?
        I hope so.
    2m
`
	rs, _ := parser.Parse(original)
	reconciler := NewReconcilerAtRecord(rs, Ɀ_Date_(2018, 1, 1))
	require.NotNil(t, reconciler)
	result, err := reconciler.CloseOpenRange(Ɀ_Time_(16, 42), "Yes!")
	require.Nil(t, err)
	assert.Equal(t, `
2018-01-01
    1h Multiline...
        ...entry summary
    15:00-16:42 Will this close?
        I hope so. Yes!
    2m
`, result.AllSerialised)
}

func TestReconcilerClosesOpenRangeWithExtendingSummaryOnNextLine(t *testing.T) {
	original := `
2018-01-01
    16:00-? Started...
    -45m break
`
	rs, _ := parser.Parse(original)
	reconciler := NewReconcilerAtRecord(rs, Ɀ_Date_(2018, 1, 1))
	require.NotNil(t, reconciler)
	result, err := reconciler.CloseOpenRange(Ɀ_Time_(18, 01), "\nStopped.")
	require.Nil(t, err)
	assert.Equal(t, `
2018-01-01
    16:00-18:01 Started...
        Stopped.
    -45m break
`, result.AllSerialised)
}
