package parser

import (
	"github.com/stretchr/testify/assert"
	"klog"
	"testing"
)

func TestSerialiseNoRecordsToEmptyString(t *testing.T) {
	text := SerialiseRecords([]src.Record{}, nil)
	assert.Equal(t, "", text)
}

func TestSerialiseEndsWithNewlineIfContainsContent(t *testing.T) {
	text := SerialiseRecord(src.NewRecord(src.Ɀ_Date_(2020, 01, 15)), nil)
	lastChar := []rune(text)[len(text)-1]
	assert.Equal(t, '\n', lastChar)
}

func TestSerialiseRecordWithMinimalRecord(t *testing.T) {
	text := SerialiseRecord(src.NewRecord(src.Ɀ_Date_(2020, 01, 15)), nil)
	assert.Equal(t, `2020-01-15
`, text)
}

func TestSerialiseRecordWithCompleteRecord(t *testing.T) {
	r := src.NewRecord(src.Ɀ_Date_(2020, 01, 15))
	r.SetShouldTotal(src.NewDuration(7, 30))
	_ = r.SetSummary("This is a\nmultiline summary")
	r.AddRange(src.Ɀ_Range_(src.Ɀ_Time_(8, 00), src.Ɀ_Time_(12, 15)), "Foo")
	r.AddDuration(src.NewDuration(2, 15), "Bar")
	_ = r.StartOpenRange(src.Ɀ_Time_(14, 38), "Baz")
	r.AddDuration(src.NewDuration(-1, -51), "")
	r.AddRange(src.Ɀ_Range_(src.Ɀ_TimeYesterday_(23, 23), src.Ɀ_Time_(4, 3)), "")
	r.AddRange(src.Ɀ_Range_(src.Ɀ_Time_(22, 0), src.Ɀ_TimeTomorrow_(0, 1)), "")
	text := SerialiseRecord(r, nil)
	assert.Equal(t, `2020-01-15 (7h30m!)
This is a
multiline summary
    8:00 - 12:15 Foo
    2h15m Bar
    14:38 - ? Baz
    -1h51m
    <23:23 - 4:03
    22:00 - 0:01>
`, text)
}

func TestSerialiseMultipleRecords(t *testing.T) {
	text := SerialiseRecords([]src.Record{
		src.NewRecord(src.Ɀ_Date_(2020, 01, 15)),
		src.NewRecord(src.Ɀ_Date_(2020, 01, 20)),
	}, nil)
	assert.Equal(t, `2020-01-15

2020-01-20
`, text)
}
