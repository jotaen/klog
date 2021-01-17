package parser

import (
	"github.com/stretchr/testify/assert"
	. "klog/record"
	"testing"
)

func TestSerialiseNoRecordsToEmptyString(t *testing.T) {
	text := SerialiseRecords([]Record{}, nil)
	assert.Equal(t, "", text)
}

func TestSerialiseEndsWithNewlineIfContainsContent(t *testing.T) {
	text := SerialiseRecord(NewRecord(Ɀ_Date_(2020, 01, 15)), nil)
	lastChar := []rune(text)[len(text)-1]
	assert.Equal(t, '\n', lastChar)
}

func TestSerialiseRecordWithMinimalRecord(t *testing.T) {
	text := SerialiseRecord(NewRecord(Ɀ_Date_(2020, 01, 15)), nil)
	assert.Equal(t, `2020-01-15
`, text)
}

func TestSerialiseRecordWithCompleteRecord(t *testing.T) {
	r := NewRecord(Ɀ_Date_(2020, 01, 15))
	r.SetShouldTotal(NewDuration(7, 30))
	_ = r.SetSummary("This is a\nmultiline summary")
	r.AddRange(Ɀ_Range_(Ɀ_Time_(8, 00), Ɀ_Time_(12, 15)), "Foo")
	r.AddDuration(NewDuration(2, 15), "Bar")
	_ = r.StartOpenRange(Ɀ_Time_(14, 38), "Baz")
	r.AddDuration(NewDuration(-1, -51), "")
	r.AddRange(Ɀ_Range_(Ɀ_TimeYesterday_(23, 23), Ɀ_Time_(4, 3)), "")
	r.AddRange(Ɀ_Range_(Ɀ_Time_(22, 0), Ɀ_TimeTomorrow_(0, 1)), "")
	text := SerialiseRecord(r, nil)
	assert.Equal(t, `2020-01-15 (7h30m!)
This is a
multiline summary
	8:00 - 12:15 Foo
	2h15m Bar
	14:38 - Baz
	-1h51m
	<23:23 - 4:03
	22:00 - 0:01>
`, text)
}

func TestSerialiseMultipleRecords(t *testing.T) {
	text := SerialiseRecords([]Record{
		NewRecord(Ɀ_Date_(2020, 01, 15)),
		NewRecord(Ɀ_Date_(2020, 01, 20)),
	}, nil)
	assert.Equal(t, `2020-01-15

2020-01-20
`, text)
}
