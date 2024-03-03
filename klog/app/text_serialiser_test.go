package app

import (
	"github.com/jotaen/klog/klog"
	tf "github.com/jotaen/klog/klog/app/cli/terminalformat"
	"github.com/jotaen/klog/klog/parser"
	"github.com/stretchr/testify/assert"
	"testing"
)

var serialiser = NewSerialiser(tf.NewStyler(tf.COLOUR_THEME_NO_COLOUR), false)

func TestSerialiseNoRecordsToEmptyString(t *testing.T) {
	text := parser.SerialiseRecords(serialiser, []klog.Record{}...).ToString()
	assert.Equal(t, "", text)
}

func TestSerialiseEndsWithNewlineIfContainsContent(t *testing.T) {
	text := parser.SerialiseRecords(serialiser, klog.NewRecord(klog.Ɀ_Date_(2020, 01, 15))).ToString()
	lastChar := []rune(text)[len(text)-1]
	assert.Equal(t, '\n', lastChar)
}

func TestSerialiseRecordWithMinimalRecord(t *testing.T) {
	text := parser.SerialiseRecords(serialiser, klog.NewRecord(klog.Ɀ_Date_(2020, 01, 15))).ToString()
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
	_ = r.Start(klog.NewOpenRange(klog.Ɀ_Time_(14, 38)), klog.Ɀ_EntrySummary_("Baz"))
	r.AddDuration(klog.NewDuration(-1, -51), nil)
	r.AddRange(klog.Ɀ_Range_(klog.Ɀ_TimeYesterday_(23, 23), klog.Ɀ_Time_(4, 3)), nil)
	r.AddRange(klog.Ɀ_Range_(klog.Ɀ_Time_(22, 0), klog.Ɀ_TimeTomorrow_(0, 1)), nil)
	text := parser.SerialiseRecords(serialiser, r).ToString()
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
	text := parser.SerialiseRecords(serialiser, []klog.Record{
		klog.NewRecord(klog.Ɀ_Date_(2020, 01, 15)),
		klog.NewRecord(klog.Ɀ_Date_(2020, 01, 20)),
	}...).ToString()
	assert.Equal(t, `2020-01-15

2020-01-20
`, text)
}

func TestParseAndSerialiseCycle(t *testing.T) {
	for _, txt := range texts {
		p := parser.NewSerialParser()
		rs, _, _ := p.Parse(txt)
		s := parser.SerialiseRecords(serialiser, rs...).ToString()
		assert.Equal(t, txt, s)
	}
}

var texts = []string{
	// Empty document.
	``,

	// Minimal document.
	`2015-05-14
`,

	// Preserves non-canonical formatting variants.
	`2015/11/28
    +1h
    2:00am-3:12pm
    12:00 - ????????????????
`,

	// Non-ASCII characters.
	`2000-01-01
日本語を母語とする大和民族が国民のほとんどを占める。自然地理的には、
ユーラシア大陸の東に位置しており、環太平洋火山帯を構成する。
島嶼国であり、領土が海に囲まれているため地続きの国境は存在しない。
日本列島は本州、北海道、九州、四国、沖縄島（以上本土）
も含めて6852の島を有する。
    1h 🙂🥸🤠👍🏽

2018-01-05
मुख्य #रूपमा काम
    10:00-12:30
        बगैचा खन्नुहोस्
    1:00am-3:00pm
         कर #घोषणा
`,

	// Longer document with all kinds of variants.
	`1999-05-31 (8h30m!)
Summary that consists of multiple
lines and contains a #tag as well.
    5h30m This and that
    -2h Something else
    +12m
    0m
    +0m
    -0m
    <18:00 - 4:00 Foo
        Bar
    19:00 - 20:00
                            Baz
                          Bar
    19:00 - 20:00
    20:01 - 0:15>
    1:00am - 3:12pm
    7:00 - ?

2000-02-12
    <18:00-4:00
    12:00-??????????

2018-01-04 (3m!)
    1h Домашня робота 🏡...
    2h Сьогодні я дзвонив
        Дімі і складав плани

2018-01-06
    +3h sázet květiny
    14:00 - ? jít na #procházku, vynést
        odpadky, #přines noviny
`,
}
