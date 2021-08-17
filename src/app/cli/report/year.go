package report

import (
	"fmt"
	. "klog"
	"klog/lib/jotaen/terminalformat"
	"klog/service"
)

type yearAggregator struct{}

func NewYearAggregator() Aggregator {
	return &yearAggregator{}
}

func (a *yearAggregator) NumberOfPrefixColumns() int {
	return 1
}

func (a *yearAggregator) DateHash(date Date) Hash {
	return Hash(service.NewYearHash(date))
}

func (a *yearAggregator) OnHeaderPrefix(table *terminalformat.Table) {
	table.
		CellL("    ") // 2020
}

func (a *yearAggregator) OnRowPrefix(table *terminalformat.Table, date Date) {
	// Year
	table.CellR(fmt.Sprint(date.Year()))
}
