package report

import (
	"fmt"
	"github.com/jotaen/klog/lib/jotaen/terminalformat"
	. "github.com/jotaen/klog/src"
	"github.com/jotaen/klog/src/service"
)

type weekAggregator struct {
	y int
}

func NewWeekAggregator() Aggregator {
	return &weekAggregator{-1}
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
		table.CellR(fmt.Sprint(date.Year()))
		a.y = date.Year()
	} else {
		table.Skip(1)
	}

	// Week
	table.CellR(fmt.Sprintf("Week %2v", date.WeekNumber()))
}
