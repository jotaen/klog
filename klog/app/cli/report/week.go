package report

import (
	"fmt"
	"github.com/jotaen/klog/klog"
	"github.com/jotaen/klog/klog/app/cli/terminalformat"
	"github.com/jotaen/klog/klog/service/period"
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

func (a *weekAggregator) DateHash(date klog.Date) period.Hash {
	return period.Hash(period.NewWeekFromDate(date).Hash())
}

func (a *weekAggregator) OnHeaderPrefix(table *terminalformat.Table) {
	table.
		CellL("    ").    // 2020
		CellL("        ") // Week 33
}

func (a *weekAggregator) OnRowPrefix(table *terminalformat.Table, date klog.Date) {
	year, week := date.WeekNumber()

	if year != a.y {
		table.CellR(fmt.Sprint(year))
		a.y = year
	} else {
		table.Skip(1)
	}

	table.CellR(fmt.Sprintf("Week %2v", week))
}
