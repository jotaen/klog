package report

import (
	"fmt"
	. "github.com/jotaen/klog/src"
	"github.com/jotaen/klog/src/app/lib/terminalformat"
	"github.com/jotaen/klog/src/service/period"
)

type quarterAggregator struct {
	y int
}

func NewQuarterAggregator() Aggregator {
	return &quarterAggregator{-1}
}

func (a *quarterAggregator) NumberOfPrefixColumns() int {
	return 2
}

func (a *quarterAggregator) DateHash(date Date) period.Hash {
	return period.Hash(period.NewQuarterFromDate(date).Hash())
}

func (a *quarterAggregator) OnHeaderPrefix(table *terminalformat.Table) {
	table.
		CellL("    "). // 2020
		CellL("  ")    // Q2
}

func (a *quarterAggregator) OnRowPrefix(table *terminalformat.Table, date Date) {
	// Year
	if date.Year() != a.y {
		table.CellR(fmt.Sprint(date.Year()))
		a.y = date.Year()
	} else {
		table.Skip(1)
	}

	// Quarter
	table.CellR(fmt.Sprintf("Q%1v", date.Quarter()))
}
