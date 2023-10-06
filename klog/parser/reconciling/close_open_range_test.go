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
	atTime := klog.Ɀ_Time_(15, 30)
	reconciler := NewReconcilerAtRecord(atDate)(rs, bs)
	require.NotNil(t, reconciler)
	err := reconciler.CloseOpenRange(atTime, NoReformat[klog.TimeFormat](), nil)
	require.Nil(t, err)

	result := assertResult(t, reconciler)
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
	atTime := klog.Ɀ_Time_(15, 22)
	reconciler := NewReconcilerAtRecord(atDate)(rs, bs)
	require.NotNil(t, reconciler)
	err := reconciler.CloseOpenRange(atTime, NoReformat[klog.TimeFormat](), klog.Ɀ_EntrySummary_("Finished."))
	require.Nil(t, err)

	result := assertResult(t, reconciler)
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
	atTime := klog.Ɀ_Time_(15, 22)
	reconciler := NewReconcilerAtRecord(atDate)(rs, bs)
	require.NotNil(t, reconciler)
	err := reconciler.CloseOpenRange(atTime, NoReformat[klog.TimeFormat](), klog.Ɀ_EntrySummary_("", "Finished."))
	require.Nil(t, err)

	result := assertResult(t, reconciler)
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
	atTime := klog.Ɀ_Time_(16, 42)
	reconciler := NewReconcilerAtRecord(atDate)(rs, bs)
	require.NotNil(t, reconciler)
	err := reconciler.CloseOpenRange(atTime, NoReformat[klog.TimeFormat](), klog.Ɀ_EntrySummary_("Yes!"))
	require.Nil(t, err)

	result := assertResult(t, reconciler)
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
	atTime := klog.Ɀ_Time_(16, 15)
	reconciler := NewReconcilerAtRecord(atDate)(rs, bs)
	require.NotNil(t, reconciler)
	err := reconciler.CloseOpenRange(atTime, NoReformat[klog.TimeFormat](), klog.Ɀ_EntrySummary_("🪴"))
	require.Nil(t, err)

	result := assertResult(t, reconciler)
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
	atTime := klog.Ɀ_Time_(18, 01)
	reconciler := NewReconcilerAtRecord(atDate)(rs, bs)
	require.NotNil(t, reconciler)
	err := reconciler.CloseOpenRange(atTime, NoReformat[klog.TimeFormat](), klog.Ɀ_EntrySummary_("", "Stopped."))
	require.Nil(t, err)

	result := assertResult(t, reconciler)
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
	atTime := klog.Ɀ_Time_(15, 30)
	reconciler := NewReconcilerAtRecord(atDate)(rs, bs)
	require.NotNil(t, reconciler)
	err := reconciler.CloseOpenRange(atTime, ReformatAutoStyle[klog.TimeFormat](), nil)
	require.Nil(t, err)

	result := assertResult(t, reconciler)
	assert.Equal(t, `
2010-04-27
    3:00pm - 3:30pm
`, result.AllSerialised)
}

func TestReconcilerClosesOpenRangeWithExplicitOverrideStyle(t *testing.T) {
	original := `
2010-04-27
    3:00pm - ??
`
	rs, bs, _ := parser.NewSerialParser().Parse(original)
	atDate := klog.Ɀ_Date_(2010, 4, 27)
	atTime := klog.Ɀ_Time_(15, 30)
	reconciler := NewReconcilerAtRecord(atDate)(rs, bs)
	require.NotNil(t, reconciler)
	err := reconciler.CloseOpenRange(atTime, ReformatExplicitly[klog.TimeFormat](klog.TimeFormat{Use24HourClock: true}), nil)
	require.Nil(t, err)

	result := assertResult(t, reconciler)
	assert.Equal(t, `
2010-04-27
    3:00pm - 15:30
`, result.AllSerialised)
}

func TestReconcilerClosesOpenRangeWithOwnStyle(t *testing.T) {
	original := `
2010-04-27
    3:00pm - ??
`
	rs, bs, _ := parser.NewSerialParser().Parse(original)
	atDate := klog.Ɀ_Date_(2010, 4, 27)
	atTime := klog.Ɀ_Time_(15, 30) // Not an am/pm time!
	reconciler := NewReconcilerAtRecord(atDate)(rs, bs)
	require.NotNil(t, reconciler)
	err := reconciler.CloseOpenRange(atTime, NoReformat[klog.TimeFormat](), nil)
	require.Nil(t, err)

	result := assertResult(t, reconciler)
	assert.Equal(t, `
2010-04-27
    3:00pm - 15:30
`, result.AllSerialised)
}
