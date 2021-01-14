package parser

import (
	"github.com/stretchr/testify/assert"
	"klog/record"
	"testing"
)

func TestSerialiseNoRecordsToEmptyString(t *testing.T) {
	text := Serialise(nil)
	assert.Equal(t, "", text)
}

func TestSerialiseEndsWithNewlineIfContainsContent(t *testing.T) {
	text := Serialise([]record.Record{
		record.NewRecord(record.Ɀ_Date_(2020, 01, 15)),
	})
	lastChar := []rune(text)[len(text)-1]
	assert.Equal(t, '\n', lastChar)
}

func TestSerialiseRecordWithMinimalRecord(t *testing.T) {
	text := Serialise([]record.Record{
		record.NewRecord(record.Ɀ_Date_(2020, 01, 15)),
	})
	assert.Equal(t, `2020-01-15
`, text)
}

func TestSerialiseRecordWithCompleteRecord(t *testing.T) {
	r := record.NewRecord(record.Ɀ_Date_(2020, 01, 15))
	r.SetSummary("This is a\nmultiline summary")
	r.AddRange(record.Ɀ_Range_(record.Ɀ_Time_(8, 00), record.Ɀ_Time_(12, 15)))
	r.AddDuration(record.NewDuration(2, 15))
	r.StartOpenRange(record.Ɀ_Time_(14, 38))
	r.AddDuration(record.NewDuration(-1, -51))
	r.AddRange(record.Ɀ_Range_(record.Ɀ_TimeYesterday_(23, 23), record.Ɀ_Time_(4, 3)))
	r.AddRange(record.Ɀ_Range_(record.Ɀ_Time_(22, 0), record.Ɀ_TimeTomorrow_(0, 1)))
	text := Serialise([]record.Record{r})
	assert.Equal(t, `2020-01-15
This is a
multiline summary
	8:00 - 12:15
	2h15m
	14:38 -
	-1h51m
	<23:23 - 4:03
	22:00 - 0:01>
`, text)
}
