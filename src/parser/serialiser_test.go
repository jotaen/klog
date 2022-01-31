package parser

import (
	. "github.com/jotaen/klog/src"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestSerialiseNoRecordsToEmptyString(t *testing.T) {
	text := PlainSerialiser.SerialiseRecords([]Record{}...)
	assert.Equal(t, "", text)
}

func TestSerialiseEndsWithNewlineIfContainsContent(t *testing.T) {
	text := PlainSerialiser.SerialiseRecords(NewRecord(Ɀ_Date_(2020, 01, 15)))
	lastChar := []rune(text)[len(text)-1]
	assert.Equal(t, '\n', lastChar)
}

func TestSerialiseRecordWithMinimalRecord(t *testing.T) {
	text := PlainSerialiser.SerialiseRecords(NewRecord(Ɀ_Date_(2020, 01, 15)))
	assert.Equal(t, `2020-01-15
`, text)
}

func TestSerialiseRecordWithCompleteRecord(t *testing.T) {
	r := NewRecord(Ɀ_Date_(2020, 01, 15))
	r.SetShouldTotal(NewDuration(7, 30))
	r.SetSummary(Ɀ_RecordSummary_("This is a", "multiline summary"))
	r.AddRange(Ɀ_Range_(Ɀ_Time_(8, 00), Ɀ_Time_(12, 15)), Ɀ_EntrySummary_("Foo"))
	r.AddDuration(NewDuration(2, 15), Ɀ_EntrySummary_("Bar", "asdf"))
	r.AddDuration(NewDuration(0, 0), Ɀ_EntrySummary_("", "Summary text...", "...more text...", "    ....preceding whitespace is ok"))
	_ = r.StartOpenRange(Ɀ_Time_(14, 38), Ɀ_EntrySummary_("Baz"))
	r.AddDuration(NewDuration(-1, -51), nil)
	r.AddRange(Ɀ_Range_(Ɀ_TimeYesterday_(23, 23), Ɀ_Time_(4, 3)), nil)
	r.AddRange(Ɀ_Range_(Ɀ_Time_(22, 0), Ɀ_TimeTomorrow_(0, 1)), nil)
	text := PlainSerialiser.SerialiseRecords(r)
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
	text := PlainSerialiser.SerialiseRecords([]Record{
		NewRecord(Ɀ_Date_(2020, 01, 15)),
		NewRecord(Ɀ_Date_(2020, 01, 20)),
	}...)
	assert.Equal(t, `2020-01-15

2020-01-20
`, text)
}
