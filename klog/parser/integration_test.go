package parser

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestParseAndSerialiseCycle(t *testing.T) {
	text := `1999-05-31 (8h30m!)
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

2018-01-05
मुख्य #रूपमा काम
    10:00-12:30 बगैचा खन्नुहोस्
    1:00am-3:00pm कर #घोषणा

2018-01-06
    +3h sázet květiny
    14:00 - ? jít na #procházku, vynést
        odpadky, #přines noviny
`
	rs, _, _ := NewSerialParser().Parse(text)
	require.Len(t, rs, 5)
	s := SerialiseRecords(PlainSerialiser{}, rs...).ToString()
	assert.Equal(t, text, s)
}
