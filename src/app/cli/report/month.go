package report

import (
	"fmt"
	. "klog"
	"klog/app/cli/lib"
	"klog/lib/jotaen/terminalformat"
	"klog/service"
)

type monthAggregator struct {
	y int
	m int
}

func NewMonthAggregator() Aggregator {
	return &monthAggregator{-1, -1}
}

func (a *monthAggregator) NumberOfPrefixColumns() int {
	return 2
}

func (a *monthAggregator) DateHash(date Date) Hash {
	return Hash(service.NewMonthHash(date))
}

func (a *monthAggregator) OnHeaderPrefix(table *terminalformat.Table) {
	table.
		CellL("    "). // 2020
		CellL("   ")   // Dec
}

func (a *monthAggregator) OnRowPrefix(table *terminalformat.Table, date Date) {
	// Year
	if date.Year() != a.y {
		a.m = -1 // force month to be recalculated
		table.CellR(fmt.Sprint(date.Year()))
		a.y = date.Year()
	} else {
		table.Skip(1)
	}

	// Month
	if date.Month() != a.m {
		a.m = date.Month()
		table.CellR(lib.PrettyMonth(a.m)[:3])
	} else {
		table.Skip(1)
	}
}
