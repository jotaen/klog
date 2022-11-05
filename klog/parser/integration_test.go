package parser

import (
	"github.com/jotaen/klog/klog"
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
`
	prs, _, _ := NewSerialParser().Parse(text)
	require.Len(t, prs, 2)
	rs := make([]klog.Record, len(prs))
	for i, pr := range prs {
		rs[i] = klog.Record(pr)
	}
	s := SerialiseRecords(PlainSerialiser{}, rs...).ToString()
	assert.Equal(t, text, s)
}
