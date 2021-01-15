package parser

import (
	"github.com/stretchr/testify/assert"
	. "klog/record"
	"testing"
)

func TestSerialiseNoRecordsToEmptyString(t *testing.T) {
	text := ToPlainText(nil)
	assert.Equal(t, "", text)
}

func TestSerialiseEndsWithNewlineIfContainsContent(t *testing.T) {
	text := ToPlainText([]Record{
		NewRecord(Ɀ_Date_(2020, 01, 15)),
	})
	lastChar := []rune(text)[len(text)-1]
	assert.Equal(t, '\n', lastChar)
}

func TestSerialiseRecordWithMinimalRecord(t *testing.T) {
	text := ToPlainText([]Record{
		NewRecord(Ɀ_Date_(2020, 01, 15)),
	})
	assert.Equal(t, `2020-01-15
`, text)
}

func TestSerialiseRecordWithCompleteRecord(t *testing.T) {
	r := NewRecord(Ɀ_Date_(2020, 01, 15))
	r.SetShouldTotal(NewDuration(7, 30))
	r.SetSummary("This is a\nmultiline summary")
	r.AddRange(Ɀ_Range_(Ɀ_Time_(8, 00), Ɀ_Time_(12, 15)))
	r.AddDuration(NewDuration(2, 15))
	r.StartOpenRange(Ɀ_Time_(14, 38))
	r.AddDuration(NewDuration(-1, -51))
	r.AddRange(Ɀ_Range_(Ɀ_TimeYesterday_(23, 23), Ɀ_Time_(4, 3)))
	r.AddRange(Ɀ_Range_(Ɀ_Time_(22, 0), Ɀ_TimeTomorrow_(0, 1)))
	text := ToPlainText([]Record{r})
	assert.Equal(t, `2020-01-15 (7h30m!)
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
