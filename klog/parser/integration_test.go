package parser

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestParseAndSerialiseCycle(t *testing.T) {
	for _, txt := range texts {
		for _, p := range parsers {
			rs, _, _ := p.Parse(txt)
			s := SerialiseRecords(PlainSerialiser{}, rs...).ToString()
			assert.Equal(t, txt, s)
		}
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
