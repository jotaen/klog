package report

import (
	"fmt"
	"github.com/jotaen/klog/klog"
	"github.com/jotaen/klog/klog/service/period"
	tf "github.com/jotaen/klog/lib/terminalformat"
)

type yearAggregator struct{}

func NewYearAggregator() Aggregator {
	return &yearAggregator{}
}

func (a *yearAggregator) NumberOfPrefixColumns() int {
	return 1
}

func (a *yearAggregator) DateHash(date klog.Date) period.Hash {
	return period.Hash(period.NewYearFromDate(date).Hash())
}

func (a *yearAggregator) OnHeaderPrefix(table *tf.Table) {
	table.
		CellL("    ") // 2020
}

func (a *yearAggregator) OnRowPrefix(table *tf.Table, date klog.Date) {
	// Year
	table.CellR(fmt.Sprint(date.Year()))
}
