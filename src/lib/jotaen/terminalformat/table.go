package terminalformat

import "strings"

type Table struct {
	cells           []string
	alignments      []Alignment
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
		cells:           []string{},
		numberOfColumns: numberOfColumns,
		longestCell:     make([]int, numberOfColumns),
		currentColumn:   0,
		columnSeparator: columnSeparator,
	}
}

func (t *Table) cell(text string, alignment Alignment) {
	t.cells = append(t.cells, text)
	t.alignments = append(t.alignments, alignment)
	if len(text) > t.longestCell[t.currentColumn] {
		t.longestCell[t.currentColumn] = len(text)
	}
	t.currentColumn++
	if t.currentColumn >= t.numberOfColumns {
		t.currentColumn = 0
	}
}

func (t *Table) ToString() string {
	result := ""
	for i, cell := range t.cells {
		col := i % t.numberOfColumns
		if i > 0 && col == 0 {
			result += "\n"
		}
		if col > 0 {
			result += t.columnSeparator
		}
		padding := strings.Repeat(" ", t.longestCell[col]-len(cell))
		if t.alignments[i] == ALIGN_RIGHT {
			result += padding
		}
		result += cell
		if t.alignments[i] == ALIGN_LEFT {
			result += padding
		}
	}
	return result
}
