package reconciling

import (
	"github.com/jotaen/klog/klog"
	"github.com/jotaen/klog/klog/parser"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestReconcilerAppendsToOpenRangeSummary(t *testing.T) {
	original := `
2018-01-01
    8:00 - 9:00 Before
    9:00 - ? First line
        Second line
    1h After
`
	rs, bs, _ := parser.NewSerialParser().Parse(original)
	reconciler := NewReconcilerAtRecord(klog.Ɀ_Date_(2018, 1, 1))(rs, bs)
	require.NotNil(t, reconciler)
	err := reconciler.AppendToOpenRange(klog.Ɀ_EntrySummary_("", "Third line #tag"))
	require.Nil(t, err)

	result := assertResult(t, reconciler)
	assert.Equal(t, `
2018-01-01
    8:00 - 9:00 Before
    9:00 - ? First line
        Second line
        Third line #tag
    1h After
`, result.AllSerialised)
}

func TestReconcilerAppendToOpenRangeFailsWithoutOpenRange(t *testing.T) {
	original := `
2018-01-01
    8:00 - 9:00
`
	rs, bs, _ := parser.NewSerialParser().Parse(original)
	reconciler := NewReconcilerAtRecord(klog.Ɀ_Date_(2018, 1, 1))(rs, bs)
	require.NotNil(t, reconciler)
	err := reconciler.AppendToOpenRange(klog.Ɀ_EntrySummary_("", "Text"))
	require.Error(t, err)
}

func TestReconcilerAppendsToNthEntrySummary(t *testing.T) {
	original := `
2018-01-01
    8:00 - 9:00 First
        continued
    2h Second
    1h Third
`
	rs, bs, _ := parser.NewSerialParser().Parse(original)
	reconciler := NewReconcilerAtRecord(klog.Ɀ_Date_(2018, 1, 1))(rs, bs)
	require.NotNil(t, reconciler)
	require.Nil(t, reconciler.AppendToNthEntry(2, klog.Ɀ_EntrySummary_("entry")))
	require.Nil(t, reconciler.AppendToNthEntry(-3, klog.Ɀ_EntrySummary_("", "more")))

	result := assertResult(t, reconciler)
	assert.Equal(t, `
2018-01-01
    8:00 - 9:00 First
        continued
        more
    2h Second entry
    1h Third
`, result.AllSerialised)
}

func TestReconcilerAppendToNthEntryFailsForMissingEntry(t *testing.T) {
	original := `
2018-01-01
    8:00 - 9:00
`
	rs, bs, _ := parser.NewSerialParser().Parse(original)
	reconciler := NewReconcilerAtRecord(klog.Ɀ_Date_(2018, 1, 1))(rs, bs)
	require.NotNil(t, reconciler)
	for _, nth := range []int{0, 2, -2} {
		require.Error(t, reconciler.AppendToNthEntry(nth, klog.Ɀ_EntrySummary_("Text")))
	}
}
