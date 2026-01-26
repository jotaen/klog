package report

import (
	"fmt"
	"github.com/jotaen/klog/klog"
	"github.com/jotaen/klog/klog/service/period"
	tf "github.com/jotaen/klog/lib/terminalformat"
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

func (a *quarterAggregator) DateHash(date klog.Date) period.Hash {
	return period.Hash(period.NewQuarterFromDate(date).Hash())
}

func (a *quarterAggregator) OnHeaderPrefix(table *tf.Table) {
	table.
		CellL("    "). // 2020
		CellL("  ")    // Q2
}

func (a *quarterAggregator) OnRowPrefix(table *tf.Table, date klog.Date) {
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
