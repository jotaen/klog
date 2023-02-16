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
	rs, bs, _ := parser.NewSerialParser().Parse(original)
	atDate := klog.Ɀ_Date_(2018, 1, 1)
	atTime := klog.Ɀ_Time_(8, 3)
	reconciler := NewReconcilerAtRecord(atDate)(rs, bs)
	require.NotNil(t, reconciler)
	result, err := reconciler.StartOpenRange(atTime, NoReformat[klog.TimeFormat](), nil)
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
	rs, bs, _ := parser.NewSerialParser().Parse(original)
	atDate := klog.Ɀ_Date_(2018, 1, 1)
	atTime := klog.Ɀ_Time_(8, 3)
	reconciler := NewReconcilerAtRecord(atDate)(rs, bs)
	require.NotNil(t, reconciler)
	result, err := reconciler.StartOpenRange(atTime, NoReformat[klog.TimeFormat](), klog.Ɀ_EntrySummary_("Test"))
	require.Nil(t, err)
	assert.Equal(t, `
2018-01-01
	5h22m
		Existing entry
	8:03 - ? Test
`, result.AllSerialised)
}

func TestReconcilerStartsOpenRangeWithUTF8Summary(t *testing.T) {
	original := `
2018-01-01
ኣብ ቤት ጽሕፈት ሓደ መዓልቲ።
	8:00 - 12:00 ንኽሰርሕ
`
	rs, bs, _ := parser.NewSerialParser().Parse(original)
	atDate := klog.Ɀ_Date_(2018, 1, 1)
	atTime := klog.Ɀ_Time_(12, 00)
	reconciler := NewReconcilerAtRecord(atDate)(rs, bs)
	require.NotNil(t, reconciler)
	result, err := reconciler.StartOpenRange(atTime, NoReformat[klog.TimeFormat](), klog.Ɀ_EntrySummary_("ናይ ምሳሕ ዕረፍቲ"))
	require.Nil(t, err)
	assert.Equal(t, `
2018-01-01
ኣብ ቤት ጽሕፈት ሓደ መዓልቲ።
	8:00 - 12:00 ንኽሰርሕ
	12:00 - ? ናይ ምሳሕ ዕረፍቲ
`, result.AllSerialised)
}

func TestReconcilerStartsOpenRangeWithNewMultilineSummary(t *testing.T) {
	original := `
2018-01-01
	5h22m
`
	rs, bs, _ := parser.NewSerialParser().Parse(original)
	atDate := klog.Ɀ_Date_(2018, 1, 1)
	atTime := klog.Ɀ_Time_(8, 3)
	reconciler := NewReconcilerAtRecord(atDate)(rs, bs)
	require.NotNil(t, reconciler)
	result, err := reconciler.StartOpenRange(atTime, NoReformat[klog.TimeFormat](), klog.Ɀ_EntrySummary_("", "Started...", "something!"))
	require.Nil(t, err)
	assert.Equal(t, `
2018-01-01
	5h22m
	8:03 - ?
		Started...
		something!
`, result.AllSerialised)
}

func TestReconcilerStartsOpenRangeWithAutoStyleFromSameRecord(t *testing.T) {
	original := `
2018-01-01
	2:00am-3:00am
`
	rs, bs, _ := parser.NewSerialParser().Parse(original)
	atDate := klog.Ɀ_Date_(2018, 1, 1)
	atTime := klog.Ɀ_Time_(8, 3)
	reconciler := NewReconcilerAtRecord(atDate)(rs, bs)
	require.NotNil(t, reconciler)
	result, err := reconciler.StartOpenRange(atTime, ReformatAutoStyle[klog.TimeFormat](), nil)
	require.Nil(t, err)
	// Conforms to both am/pm and spaces around dash
	assert.Equal(t, `
2018-01-01
	2:00am-3:00am
	8:03am-?
`, result.AllSerialised)
}

func TestReconcilerStartsOpenRangeWithAutoStyleFromOtherRecord(t *testing.T) {
	original := `
2018/01/01
  2:00am-?????????????????
`
	rs, bs, _ := parser.NewSerialParser().Parse(original)
	atDate := klog.Ɀ_Date_(2018, 1, 2)
	atTime := klog.Ɀ_Time_(8, 3)
	reconciler := NewReconcilerForNewRecord(atDate, ReformatAutoStyle[klog.DateFormat](), AdditionalData{})(rs, bs)
	require.NotNil(t, reconciler)
	result, err := reconciler.StartOpenRange(atTime, ReformatAutoStyle[klog.TimeFormat](), nil)
	require.Nil(t, err)
	// Conforms to both am/pm and spaces around dash
	assert.Equal(t, `
2018/01/01
  2:00am-?????????????????

2018/01/02
  8:03am-?????????????????
`, result.AllSerialised)
}

func TestReconcilerStartsOpenRangeExplicitFormatOverrulesAutoFormat(t *testing.T) {
	original := `
2018-01-01
`
	rs, bs, _ := parser.NewSerialParser().Parse(original)
	atDate := klog.Ɀ_Slashes_(klog.Ɀ_Date_(2018, 1, 2))
	atTime := klog.Ɀ_IsAmPm_(klog.Ɀ_Time_(8, 3))
	reconciler := NewReconcilerForNewRecord(atDate, ReformatExplicitly[klog.DateFormat](atDate.Format()), AdditionalData{})(rs, bs)
	require.NotNil(t, reconciler)
	result, err := reconciler.StartOpenRange(atTime, ReformatExplicitly[klog.TimeFormat](atTime.Format()), nil)
	require.Nil(t, err)
	assert.Equal(t, `
2018-01-01

2018/01/02
    8:03am - ?
`, result.AllSerialised)
}
