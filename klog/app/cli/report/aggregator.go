/*
Package report is a utility for the report command.
*/
package report

import (
	"github.com/jotaen/klog/klog"
	"github.com/jotaen/klog/klog/app/cli/terminalformat"
	"github.com/jotaen/klog/klog/service/period"
)

type Aggregator interface {
	NumberOfPrefixColumns() int
	DateHash(klog.Date) period.Hash
	OnHeaderPrefix(*terminalformat.Table)
	OnRowPrefix(*terminalformat.Table, klog.Date)
}
