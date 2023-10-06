package reconciling

import (
	"github.com/jotaen/klog/klog"
	"github.com/jotaen/klog/klog/parser"
	"github.com/jotaen/klog/klog/service"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func assertResult(t *testing.T, r *Reconciler) *Result {
	result, err := r.MakeResult()
	require.Nil(t, err)
	return result
}

func TestReconcilerAddsNewlyCreatedEntryAtEndOfFile(t *testing.T) {
	original := "\n2018-01-01\n    1h"
	rs, bs, _ := parser.NewSerialParser().Parse(original)
	reconciler := NewReconcilerAtRecord(klog.Ɀ_Date_(2018, 1, 1))(rs, bs)
	require.NotNil(t, reconciler)
	err := reconciler.AppendEntry(klog.Ɀ_EntrySummary_("16:00-17:00"))
	require.Nil(t, err)

	result := assertResult(t, reconciler)
	assert.Equal(t, `
2018-01-01
    1h
    16:00-17:00
`, result.AllSerialised)
}

func TestReconcilerAddsNewEntryInTheMiddleOfFile(t *testing.T) {
	original := `
2018-01-01
    1h

2018-01-02
Hello World
    1h
    1h45m Multiline...
        ....entry summary

2018-01-03
    5h
`
	rs, bs, _ := parser.NewSerialParser().Parse(original)
	reconciler := NewReconcilerAtRecord(klog.Ɀ_Date_(2018, 1, 2))(rs, bs)
	require.NotNil(t, reconciler)
	err := reconciler.AppendEntry(klog.Ɀ_EntrySummary_("2h30m"))
	require.Nil(t, err)

	result := assertResult(t, reconciler)
	require.NotNil(t, result)
	require.Equal(t, 150, result.Record.Entries()[2].Duration().InMinutes())
	require.Equal(t, 315, service.Total(result.Record).InMinutes())
	assert.Equal(t, `
2018-01-01
    1h

2018-01-02
Hello World
    1h
    1h45m Multiline...
        ....entry summary
    2h30m

2018-01-03
    5h
`, result.AllSerialised)
}

func TestReconcilerCanHandleUTF8Input(t *testing.T) {
	original := "\n2018-01-01\n家事\n    1h ランドリー"
	rs, bs, _ := parser.NewSerialParser().Parse(original)

	reconciler := NewReconcilerAtRecord(klog.Ɀ_Date_(2018, 1, 1))(rs, bs)
	require.NotNil(t, reconciler)
	err := reconciler.AppendEntry(klog.Ɀ_EntrySummary_("20m 掃除機"))
	require.Nil(t, err)

	result := assertResult(t, reconciler)
	assert.Equal(t, `
2018-01-01
家事
    1h ランドリー
    20m 掃除機
`, result.AllSerialised)
}

func TestReconcilerSplitsUpSummaryText(t *testing.T) {
	original := "\n2018-01-01\n    1h"
	rs, bs, _ := parser.NewSerialParser().Parse(original)

	reconciler := NewReconcilerAtRecord(klog.Ɀ_Date_(2018, 1, 1))(rs, bs)
	require.NotNil(t, reconciler)
	err := reconciler.AppendEntry(klog.Ɀ_EntrySummary_("2h This is a", "multiline summary"))
	require.Nil(t, err)

	result := assertResult(t, reconciler)
	assert.Equal(t, `
2018-01-01
    1h
    2h This is a
        multiline summary
`, result.AllSerialised)
}

func TestReconcilerStartsSummaryTextOnNextLine(t *testing.T) {
	original := "\n2018-01-01\n    1h"
	rs, bs, _ := parser.NewSerialParser().Parse(original)

	reconciler := NewReconcilerAtRecord(klog.Ɀ_Date_(2018, 1, 1))(rs, bs)
	require.NotNil(t, reconciler)
	err := reconciler.AppendEntry(klog.Ɀ_EntrySummary_("2h", "Some activity"))
	require.Nil(t, err)

	result := assertResult(t, reconciler)
	assert.Equal(t, `
2018-01-01
    1h
    2h
        Some activity
`, result.AllSerialised)
}

func TestReconcilerRejectsInvalidEntry(t *testing.T) {
	original := "2018-01-01\n"
	rs, bs, _ := parser.NewSerialParser().Parse(original)
	reconciler := NewReconcilerAtRecord(klog.Ɀ_Date_(2018, 1, 1))(rs, bs)
	require.NotNil(t, reconciler)

	err := reconciler.AppendEntry(klog.Ɀ_EntrySummary_("this is not valid entry text"))
	assert.Nil(t, err) // This doesn’t produce an error yet, but calling `MakingResult` will!

	result, rErr := reconciler.MakeResult()
	assert.Error(t, rErr)
	assert.Nil(t, result)
}
