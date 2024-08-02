package terminalformat

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestPrintTable(t *testing.T) {
	result := ""
	styler := NewStyler(COLOUR_THEME_DARK)
	table := NewTable(3, " ")
	table.
		Cell("FIRST", Options{align: ALIGN_LEFT}).
		Cell("SECOND", Options{align: ALIGN_RIGHT}).
		Cell("THIRD", Options{align: ALIGN_RIGHT}).
		CellL("1").
		CellR("2").
		CellR("3").
		Cell("long-text", Options{align: ALIGN_LEFT}).
		Cell(styler.Props(StyleProps{IsUnderlined: true}).Format("asdf"), Options{align: ALIGN_RIGHT}).
		Fill("-").
		Skip(2).
		Cell("foo", Options{align: ALIGN_LEFT})
	table.Collect(func(x string) { result += x })
	assert.Equal(t, `FIRST     SECOND THIRD
1              2     3
long-text   `+"\x1b[0m\x1b[4m"+`asdf`+"\x1b[0m"+` -----
                 foo  
`, result)
}

func TestPrintTableWithUnicode(t *testing.T) {
	result := ""
	table := NewTable(3, " ")
	table.
		Cell("FIRST", Options{align: ALIGN_LEFT}).
		Cell("SECOND", Options{align: ALIGN_LEFT}).
		Cell("THIRD", Options{align: ALIGN_LEFT}).
		CellL("first").
		CellR("șëčøñd").
		CellR("third")
	table.Collect(func(x string) { result += x })
	assert.Equal(t, `FIRST SECOND THIRD
first șëčøñd third
`, result)
}
