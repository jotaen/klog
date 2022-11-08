package reconciling

import (
	"github.com/jotaen/klog/klog"
	"github.com/jotaen/klog/klog/parser"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestReconcilerClosesOpenRange(t *testing.T) {
	original := `
2010-04-27
    15:00 - ??
`
	rs, bs, _ := parser.NewSerialParser().Parse(original)
	atDate := klog.Ɀ_Date_(2010, 4, 27)
	atTime := Styled[klog.Time]{klog.Ɀ_Time_(15, 30), false}
	reconciler := NewReconcilerAtRecord(atDate)(rs, bs)
	require.NotNil(t, reconciler)
	result, err := reconciler.CloseOpenRange(atTime, nil)
	require.Nil(t, err)
	assert.Equal(t, `
2010-04-27
    15:00 - 15:30
`, result.AllSerialised)
}

func TestReconcilerClosesOpenRangeWithNewSummary(t *testing.T) {
	original := `
2018-01-01
    15:00 - ?
`
	rs, bs, _ := parser.NewSerialParser().Parse(original)
	atDate := klog.Ɀ_Date_(2018, 1, 1)
	atTime := Styled[klog.Time]{klog.Ɀ_Time_(15, 22), false}
	reconciler := NewReconcilerAtRecord(atDate)(rs, bs)
	require.NotNil(t, reconciler)
	result, err := reconciler.CloseOpenRange(atTime, klog.Ɀ_EntrySummary_("Finished."))
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
	rs, bs, _ := parser.NewSerialParser().Parse(original)
	atDate := klog.Ɀ_Date_(2018, 1, 1)
	atTime := Styled[klog.Time]{klog.Ɀ_Time_(15, 22), false}
	reconciler := NewReconcilerAtRecord(atDate)(rs, bs)
	require.NotNil(t, reconciler)
	result, err := reconciler.CloseOpenRange(atTime, klog.Ɀ_EntrySummary_("", "Finished."))
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
	rs, bs, _ := parser.NewSerialParser().Parse(original)
	atDate := klog.Ɀ_Date_(2018, 1, 1)
	atTime := Styled[klog.Time]{klog.Ɀ_Time_(16, 42), false}
	reconciler := NewReconcilerAtRecord(atDate)(rs, bs)
	require.NotNil(t, reconciler)
	result, err := reconciler.CloseOpenRange(atTime, klog.Ɀ_EntrySummary_("Yes!"))
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

func TestReconcilerClosesOpenRangeWithUTF8Summary(t *testing.T) {
	original := `
2018-01-01
Arbeiten rund um’s Haus… 🏡
    15:00 - ? Blümchen 🌼 planzen
`
	rs, bs, _ := parser.NewSerialParser().Parse(original)
	atDate := klog.Ɀ_Date_(2018, 1, 1)
	atTime := Styled[klog.Time]{klog.Ɀ_Time_(16, 15), false}
	reconciler := NewReconcilerAtRecord(atDate)(rs, bs)
	require.NotNil(t, reconciler)
	result, err := reconciler.CloseOpenRange(atTime, klog.Ɀ_EntrySummary_("🪴"))
	require.Nil(t, err)
	assert.Equal(t, `
2018-01-01
Arbeiten rund um’s Haus… 🏡
    15:00 - 16:15 Blümchen 🌼 planzen 🪴
`, result.AllSerialised)
}

func TestReconcilerClosesOpenRangeWithExtendingSummaryOnNextLine(t *testing.T) {
	original := `
2018-01-01
    16:00-? Started...
    -45m break
`
	rs, bs, _ := parser.NewSerialParser().Parse(original)
	atDate := klog.Ɀ_Date_(2018, 1, 1)
	atTime := Styled[klog.Time]{klog.Ɀ_Time_(18, 01), false}
	reconciler := NewReconcilerAtRecord(atDate)(rs, bs)
	require.NotNil(t, reconciler)
	result, err := reconciler.CloseOpenRange(atTime, klog.Ɀ_EntrySummary_("", "Stopped."))
	require.Nil(t, err)
	assert.Equal(t, `
2018-01-01
    16:00-18:01 Started...
        Stopped.
    -45m break
`, result.AllSerialised)
}

func TestReconcilerClosesOpenRangeDetectsStyle(t *testing.T) {
	original := `
2010-04-27
    3:00pm - ??
`
	rs, bs, _ := parser.NewSerialParser().Parse(original)
	atDate := klog.Ɀ_Date_(2010, 4, 27)
	atTime := Styled[klog.Time]{klog.Ɀ_Time_(15, 30), true}
	reconciler := NewReconcilerAtRecord(atDate)(rs, bs)
	require.NotNil(t, reconciler)
	result, err := reconciler.CloseOpenRange(atTime, nil)
	require.Nil(t, err)
	assert.Equal(t, `
2010-04-27
    3:00pm - 3:30pm
`, result.AllSerialised)
}

func TestReconcilerClosesOpenRangeWithExplicitStyle(t *testing.T) {
	original := `
2010-04-27
    3:00pm - ??
`
	rs, bs, _ := parser.NewSerialParser().Parse(original)
	atDate := klog.Ɀ_Date_(2010, 4, 27)
	atTime := Styled[klog.Time]{klog.Ɀ_Time_(15, 30), false} // Not an am/pm time!
	reconciler := NewReconcilerAtRecord(atDate)(rs, bs)
	require.NotNil(t, reconciler)
	result, err := reconciler.CloseOpenRange(atTime, nil)
	require.Nil(t, err)
	assert.Equal(t, `
2010-04-27
    3:00pm - 15:30
`, result.AllSerialised)
}
