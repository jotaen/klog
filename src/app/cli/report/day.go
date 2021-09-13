package report

import (
	"fmt"
	. "github.com/jotaen/klog/src"
	"github.com/jotaen/klog/src/app/cli/lib"
	"github.com/jotaen/klog/src/lib/jotaen/terminalformat"
	"github.com/jotaen/klog/src/service"
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

func (a *dayAggregator) DateHash(date Date) Hash {
	return Hash(service.NewDayHash(date))
}

func (a *dayAggregator) OnHeaderPrefix(table *terminalformat.Table) {
	table.
		CellL("    ").   // 2020
		CellL("   ").    // Dec
		CellL("      "). // Sun
		CellR("   ")     // 17.
}

func (a *dayAggregator) OnRowPrefix(table *terminalformat.Table, date Date) {
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

	// Day
	table.CellR(lib.PrettyDay(date.Weekday())[:3]).CellR(fmt.Sprintf("%2v.", date.Day()))
}
