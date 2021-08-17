package report

import (
	"fmt"
	. "klog"
	"klog/lib/jotaen/terminalformat"
	"klog/service"
)

type quarterAggregator struct {
	y int
	q int
}

func NewQuarterAggregator() Aggregator {
	return &quarterAggregator{-1, -1}
}

func (a *quarterAggregator) NumberOfPrefixColumns() int {
	return 2
}

func (a *quarterAggregator) DateHash(date Date) Hash {
	return Hash(service.NewQuarterHash(date))
}

func (a *quarterAggregator) OnHeaderPrefix(table *terminalformat.Table) {
	table.
		CellL("    "). // 2020
		CellL("  ")    // Q2
}

func (a *quarterAggregator) OnRowPrefix(table *terminalformat.Table, date Date) {
	// Year
	if date.Year() != a.y {
		a.q = -1 // force quarter to be recalculated
		table.CellR(fmt.Sprint(date.Year()))
		a.y = date.Year()
	} else {
		table.Skip(1)
	}

	// Quarter
	if date.Quarter() != a.q {
		table.CellR(fmt.Sprintf("Q%1v", date.Quarter()))
	} else {
		table.Skip(1)
	}
}
