package reconciling

import (
	. "github.com/jotaen/klog/src"
	"github.com/jotaen/klog/src/parser"
	"github.com/jotaen/klog/src/service"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestReconcilerAddsNewEntry(t *testing.T) {
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
	rs, _ := parser.Parse(original)
	reconciler := NewReconcilerAtRecord(rs, Ɀ_Date_(2018, 1, 2))
	require.NotNil(t, reconciler)
	result, err := reconciler.AppendEntry("2h30m")
	require.Nil(t, err)
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

func TestReconcilerAddsNewlyCreatedEntryAtEndOfFile(t *testing.T) {
	original := "\n2018-01-01\n    1h"
	rs, _ := parser.Parse(original)

	reconciler := NewReconcilerAtRecord(rs, Ɀ_Date_(2018, 1, 1))
	require.NotNil(t, reconciler)
	result, err := reconciler.AppendEntry("16:00-17:00")
	require.Nil(t, err)
	assert.Equal(t, `
2018-01-01
    1h
    16:00-17:00
`, result.AllSerialised)
}

func TestReconcilerSplitsUpSummaryText(t *testing.T) {
	original := "\n2018-01-01\n    1h"
	rs, _ := parser.Parse(original)

	reconciler := NewReconcilerAtRecord(rs, Ɀ_Date_(2018, 1, 1))
	require.NotNil(t, reconciler)
	result, err := reconciler.AppendEntry("2h This is a\nmultiline summary")
	require.Nil(t, err)
	assert.Equal(t, `
2018-01-01
    1h
    2h This is a
        multiline summary
`, result.AllSerialised)
}

func TestReconcilerStartsSummaryTextOnNextLine(t *testing.T) {
	original := "\n2018-01-01\n    1h"
	rs, _ := parser.Parse(original)

	reconciler := NewReconcilerAtRecord(rs, Ɀ_Date_(2018, 1, 1))
	require.NotNil(t, reconciler)
	result, err := reconciler.AppendEntry("2h\nSome activity")
	require.Nil(t, err)
	assert.Equal(t, `
2018-01-01
    1h
    2h
        Some activity
`, result.AllSerialised)
}

func TestReconcilerRejectsInvalidEntry(t *testing.T) {
	original := "2018-01-01\n"
	rs, _ := parser.Parse(original)
	reconciler := NewReconcilerAtRecord(rs, Ɀ_Date_(2018, 1, 1))
	require.NotNil(t, reconciler)
	result, err := reconciler.AppendEntry("this is not valid entry text")
	require.Nil(t, result)
	assert.Error(t, err)
}
