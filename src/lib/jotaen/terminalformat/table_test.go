package terminalformat

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestPrintTable(t *testing.T) {
	table := NewTable(3, " ")
	table.cell("FIRST", ALIGN_LEFT)
	table.cell("SECOND", ALIGN_RIGHT)
	table.cell("THIRD", ALIGN_RIGHT)
	table.cell("1", ALIGN_LEFT)
	table.cell("2", ALIGN_RIGHT)
	table.cell("3", ALIGN_RIGHT)
	table.cell("long-text", ALIGN_LEFT)
	table.cell("asdf", ALIGN_RIGHT)
	table.cell("", ALIGN_RIGHT)
	assert.Equal(t, `FIRST     SECOND THIRD
1              2     3
long-text   asdf      `, table.ToString())
}
