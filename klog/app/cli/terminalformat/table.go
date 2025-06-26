package terminalformat

import (
	"strings"
	"unicode/utf8"
)

type Options struct {
	fill  bool
	align Alignment
}

type cell struct {
	Options
	value string
	len   int
}

type Table struct {
	cells           []cell
	numberOfColumns int
	longestCell     []int
	currentColumn   int
	columnSeparator string
}

type Alignment int

const (
	ALIGN_LEFT Alignment = iota
	ALIGN_RIGHT
)

func NewTable(numberOfColumns int, columnSeparator string) *Table {
	if numberOfColumns <= 1 {
		panic("Column count must be greater than 1")
	}
	return &Table{
		cells:           []cell{},
		numberOfColumns: numberOfColumns,
		longestCell:     make([]int, numberOfColumns),
		currentColumn:   0,
		columnSeparator: columnSeparator,
	}
}

func (t *Table) Cell(text string, opts Options) *Table {
	c := cell{
		Options: opts,
		value:   text,
		len:     utf8.RuneCountInString(StripAllAnsiSequences(text)),
	}
	t.cells = append(t.cells, c)
	if c.len > t.longestCell[t.currentColumn] {
		t.longestCell[t.currentColumn] = c.len
	}
	t.currentColumn++
	if t.currentColumn >= t.numberOfColumns {
		t.currentColumn = 0
	}
	return t
}

func (t *Table) CellL(text string) *Table {
	return t.Cell(text, Options{align: ALIGN_LEFT})
}

func (t *Table) CellR(text string) *Table {
	return t.Cell(text, Options{align: ALIGN_RIGHT})
}

func (t *Table) Skip(numberOfCells int) *Table {
	for i := 0; i < numberOfCells; i++ {
		t.Cell("", Options{})
	}
	return t
}

func (t *Table) Fill(sequence string) *Table {
	t.Cell(sequence, Options{fill: true})
	return t
}

func (t *Table) Collect(fn func(string)) {
	for i, c := range t.cells {
		col := i % t.numberOfColumns
		if i > 0 && col == 0 {
			fn("\n")
		}
		if col > 0 {
			fn(t.columnSeparator)
		}
		if c.fill {
			fn(strings.Repeat(c.value, t.longestCell[col]))
		} else {
			padding := strings.Repeat(" ", t.longestCell[col]-c.len)
			if c.align == ALIGN_RIGHT {
				fn(padding)
			}
			fn(c.value)
			if c.align == ALIGN_LEFT {
				fn(padding)
			}
		}
	}
	if len(t.cells) > 0 {
		fn("\n")
	}
}
