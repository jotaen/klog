package report

import (
	"fmt"
	"github.com/jotaen/klog/klog"
	"github.com/jotaen/klog/klog/app/cli/terminalformat"
	"github.com/jotaen/klog/klog/app/cli/util"
	"github.com/jotaen/klog/klog/service/period"
)

type dayAggregator struct {
	y int
	m int
}

func NewDayAggregator() Aggregator {
	return &dayAggregator{-1, -1}
}

func (a *dayAggregator) NumberOfPrefixColumns() int {
	return 4
}

func (a *dayAggregator) DateHash(date klog.Date) period.Hash {
	return period.Hash(period.NewDayFromDate(date).Hash())
}

func (a *dayAggregator) OnHeaderPrefix(table *terminalformat.Table) {
	table.
		CellL("    ").   // 2020
		CellL("   ").    // Dec
		CellL("      "). // Sun
		CellR("   ")     // 17.
}

func (a *dayAggregator) OnRowPrefix(table *terminalformat.Table, date klog.Date) {
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
		table.CellR(util.PrettyMonth(a.m)[:3])
	} else {
		table.Skip(1)
	}

	// Day
	table.CellR(util.PrettyDay(date.Weekday())[:3]).CellR(fmt.Sprintf("%2v.", date.Day()))
}
