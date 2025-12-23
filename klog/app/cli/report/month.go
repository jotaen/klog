package report

import (
	"fmt"

	"github.com/jotaen/klog/klog"
	"github.com/jotaen/klog/klog/app/cli/prettify"
	"github.com/jotaen/klog/klog/service/period"
	tf "github.com/jotaen/klog/lib/terminalformat"
)

type monthAggregator struct {
	y int
}

func NewMonthAggregator() Aggregator {
	return &monthAggregator{-1}
}

func (a *monthAggregator) NumberOfPrefixColumns() int {
	return 2
}

func (a *monthAggregator) DateHash(date klog.Date) period.Hash {
	return period.Hash(period.NewMonthFromDate(date).Hash())
}

func (a *monthAggregator) OnHeaderPrefix(table *tf.Table) {
	table.
		CellL("    "). // 2020
		CellL("   ")   // Dec
}

func (a *monthAggregator) OnRowPrefix(table *tf.Table, date klog.Date) {
	// Year
	if date.Year() != a.y {
		table.CellR(fmt.Sprint(date.Year()))
		a.y = date.Year()
	} else {
		table.Skip(1)
	}

	// Month
	table.CellR(prettify.PrettyMonth(date.Month())[:3])
}
