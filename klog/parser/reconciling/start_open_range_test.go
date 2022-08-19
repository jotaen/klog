package reconciling

import (
	"github.com/jotaen/klog/klog"
	"github.com/jotaen/klog/klog/parser"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestReconcilerStartsOpenRange(t *testing.T) {
	original := `
2018-01-01
	5h22m
`
	rs, _ := parser.Parse(original)
	reconciler := NewReconcilerAtRecord(rs, klog.Ɀ_Date_(2018, 1, 1))
	require.NotNil(t, reconciler)
	result, err := reconciler.StartOpenRange(klog.Ɀ_Time_(8, 3), nil)
	require.Nil(t, err)
	assert.Equal(t, `
2018-01-01
	5h22m
	8:03 - ?
`, result.AllSerialised)
}

func TestReconcilerStartsOpenRangeWithSummary(t *testing.T) {
	original := `
2018-01-01
	5h22m
		Existing entry
`
	rs, _ := parser.Parse(original)
	reconciler := NewReconcilerAtRecord(rs, klog.Ɀ_Date_(2018, 1, 1))
	require.NotNil(t, reconciler)
	result, err := reconciler.StartOpenRange(klog.Ɀ_Time_(8, 3), klog.Ɀ_EntrySummary_("Test"))
	require.Nil(t, err)
	assert.Equal(t, `
2018-01-01
	5h22m
		Existing entry
	8:03 - ? Test
`, result.AllSerialised)
}

func TestReconcilerStartsOpenRangeWithNewMultilineSummary(t *testing.T) {
	original := `
2018-01-01
	5h22m
`
	rs, _ := parser.Parse(original)
	reconciler := NewReconcilerAtRecord(rs, klog.Ɀ_Date_(2018, 1, 1))
	require.NotNil(t, reconciler)
	result, err := reconciler.StartOpenRange(klog.Ɀ_Time_(8, 3), klog.Ɀ_EntrySummary_("", "Started...", "something!"))
	require.Nil(t, err)
	assert.Equal(t, `
2018-01-01
	5h22m
	8:03 - ?
		Started...
		something!
`, result.AllSerialised)
}

func TestReconcilerStartsOpenRangeWithStyle(t *testing.T) {
	original := `
2018-01-01
	2:00am-3:00am
`
	rs, _ := parser.Parse(original)
	reconciler := NewReconcilerAtRecord(rs, klog.Ɀ_Date_(2018, 1, 1))
	require.NotNil(t, reconciler)
	result, err := reconciler.StartOpenRange(klog.Ɀ_Time_(8, 3), nil)
	require.Nil(t, err)
	// Conforms to both am/pm and spaces around dash
	assert.Equal(t, `
2018-01-01
	2:00am-3:00am
	8:03am-?
`, result.AllSerialised)
}

func TestReconcilerStartsOpenRangeWithStyleFromOtherRecord(t *testing.T) {
	original := `
2018-01-01
  2:00am-3:00am

2018-01-02
`
	rs, _ := parser.Parse(original)
	reconciler := NewReconcilerAtRecord(rs, klog.Ɀ_Date_(2018, 1, 2))
	require.NotNil(t, reconciler)
	result, err := reconciler.StartOpenRange(klog.Ɀ_Time_(8, 3), nil)
	require.Nil(t, err)
	// Conforms to both am/pm and spaces around dash
	assert.Equal(t, `
2018-01-01
  2:00am-3:00am

2018-01-02
  8:03am-?
`, result.AllSerialised)
}
