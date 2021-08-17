package report

import (
	"fmt"
	. "klog"
	"klog/lib/jotaen/terminalformat"
	"klog/service"
)

type weekAggregator struct {
	y int
	w int
}

func NewWeekAggregator() Aggregator {
	return &weekAggregator{-1, -1}
}

func (a *weekAggregator) NumberOfPrefixColumns() int {
	return 2
}

func (a *weekAggregator) DateHash(date Date) Hash {
	return Hash(service.NewWeekHash(date))
}

func (a *weekAggregator) OnHeaderPrefix(table *terminalformat.Table) {
	table.
		CellL("    ").    // 2020
		CellL("        ") // Week 33
}

func (a *weekAggregator) OnRowPrefix(table *terminalformat.Table, date Date) {
	// Year
	if date.Year() != a.y {
		a.w = -1 // force week to be recalculated
		table.CellR(fmt.Sprint(date.Year()))
		a.y = date.Year()
	} else {
		table.Skip(1)
	}

	// Week
	if date.WeekNumber() != a.w {
		table.CellR(fmt.Sprintf("Week %2v", date.WeekNumber()))
	} else {
		table.Skip(1)
	}
}
