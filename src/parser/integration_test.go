package parser

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestParseAndSerialiseCycle(t *testing.T) {
	text := `1999-05-31 (8h30m!)
Summary that consists of multiple
lines and contains a #tag as well.
    5h30m This and that
    -2h Something else
    <18:00 - 4:00 Foo
    19:00 - 20:00
    20:01 - 0:15>
    1:00am - 3:12pm
    7:00 - ?
`
	rs, _ := Parse(text)
	s := SerialiseRecords(rs, defaultHooks{})
	assert.Equal(t, text, s)
}
