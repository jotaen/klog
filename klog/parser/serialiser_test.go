package parser

import (
	"github.com/jotaen/klog/klog"
	"github.com/stretchr/testify/assert"
	"testing"
)

type plainSerialiser struct{}

func (ps plainSerialiser) Date(x klog.Date) string               { return x.ToString() }
func (ps plainSerialiser) ShouldTotal(x klog.Duration) string    { return x.ToString() }
func (ps plainSerialiser) Summary(x SummaryText) string          { return x.ToString() }
func (ps plainSerialiser) Range(x klog.Range) string             { return x.ToString() }
func (ps plainSerialiser) OpenRange(x klog.OpenRange) string     { return x.ToString() }
func (ps plainSerialiser) Duration(x klog.Duration) string       { return x.ToString() }
func (ps plainSerialiser) SignedDuration(x klog.Duration) string { return x.ToString() }
func (ps plainSerialiser) Time(x klog.Time) string               { return x.ToString() }

func TestSerialiseNoRecordsToEmptyString(t *testing.T) {
	text := SerialiseRecords(plainSerialiser{}, []klog.Record{}...)
	assert.Equal(t, "", text)
}

func TestSerialiseEndsWithNewlineIfContainsContent(t *testing.T) {
	text := SerialiseRecords(plainSerialiser{}, klog.NewRecord(klog.Ɀ_Date_(2020, 01, 15)))
	lastChar := []rune(text)[len(text)-1]
	assert.Equal(t, '\n', lastChar)
}

func TestSerialiseRecordWithMinimalRecord(t *testing.T) {
	text := SerialiseRecords(plainSerialiser{}, klog.NewRecord(klog.Ɀ_Date_(2020, 01, 15)))
	assert.Equal(t, `2020-01-15
`, text)
}

func TestSerialiseRecordWithCompleteRecord(t *testing.T) {
	r := klog.NewRecord(klog.Ɀ_Date_(2020, 01, 15))
	r.SetShouldTotal(klog.NewDuration(7, 30))
	r.SetSummary(klog.Ɀ_RecordSummary_("This is a", "multiline summary"))
	r.AddRange(klog.Ɀ_Range_(klog.Ɀ_Time_(8, 00), klog.Ɀ_Time_(12, 15)), klog.Ɀ_EntrySummary_("Foo"))
	r.AddDuration(klog.NewDuration(2, 15), klog.Ɀ_EntrySummary_("Bar", "asdf"))
	r.AddDuration(klog.NewDuration(0, 0), klog.Ɀ_EntrySummary_("", "Summary text...", "...more text...", "    ....preceding whitespace is ok"))
	_ = r.StartOpenRange(klog.Ɀ_Time_(14, 38), klog.Ɀ_EntrySummary_("Baz"))
	r.AddDuration(klog.NewDuration(-1, -51), nil)
	r.AddRange(klog.Ɀ_Range_(klog.Ɀ_TimeYesterday_(23, 23), klog.Ɀ_Time_(4, 3)), nil)
	r.AddRange(klog.Ɀ_Range_(klog.Ɀ_Time_(22, 0), klog.Ɀ_TimeTomorrow_(0, 1)), nil)
	text := SerialiseRecords(plainSerialiser{}, r)
	assert.Equal(t, `2020-01-15 (7h30m!)
This is a
multiline summary
    8:00 - 12:15 Foo
    2h15m Bar
        asdf
    0m
        Summary text...
        ...more text...
            ....preceding whitespace is ok
    14:38 - ? Baz
    -1h51m
    <23:23 - 4:03
    22:00 - 0:01>
`, text)
}

func TestSerialiseMultipleRecords(t *testing.T) {
	text := SerialiseRecords(plainSerialiser{}, []klog.Record{
		klog.NewRecord(klog.Ɀ_Date_(2020, 01, 15)),
		klog.NewRecord(klog.Ɀ_Date_(2020, 01, 20)),
	}...)
	assert.Equal(t, `2020-01-15

2020-01-20
`, text)
}
