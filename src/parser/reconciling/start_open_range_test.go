package reconciling

import (
	. "github.com/jotaen/klog/src"
	"github.com/jotaen/klog/src/parser"
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
	reconciler := NewReconcilerAtRecord(rs, Ɀ_Date_(2018, 1, 1))
	require.NotNil(t, reconciler)
	result, err := reconciler.StartOpenRange(Ɀ_Time_(8, 3), "")
	require.Nil(t, err)
	assert.Equal(t, `
2018-01-01
	5h22m
	8:03 - ?
`, result.AllSerialised)
}

func TestReconcilerStartsOpenRangeWithNewSummary(t *testing.T) {
	original := `
2018-01-01
	5h22m
`
	rs, _ := parser.Parse(original)
	reconciler := NewReconcilerAtRecord(rs, Ɀ_Date_(2018, 1, 1))
	require.NotNil(t, reconciler)
	result, err := reconciler.StartOpenRange(Ɀ_Time_(8, 3), "Started!")
	require.Nil(t, err)
	assert.Equal(t, `
2018-01-01
	5h22m
	8:03 - ? Started!
`, result.AllSerialised)
}

func TestReconcilerStartsOpenRangeWithStyle(t *testing.T) {
	original := `
2018-01-01
	2:00am-3:00am
`
	rs, _ := parser.Parse(original)
	reconciler := NewReconcilerAtRecord(rs, Ɀ_Date_(2018, 1, 1))
	require.NotNil(t, reconciler)
	result, err := reconciler.StartOpenRange(Ɀ_Time_(8, 3), "")
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
	reconciler := NewReconcilerAtRecord(rs, Ɀ_Date_(2018, 1, 2))
	require.NotNil(t, reconciler)
	result, err := reconciler.StartOpenRange(Ɀ_Time_(8, 3), "")
	require.Nil(t, err)
	// Conforms to both am/pm and spaces around dash
	assert.Equal(t, `
2018-01-01
  2:00am-3:00am

2018-01-02
  8:03am-?
`, result.AllSerialised)
}
