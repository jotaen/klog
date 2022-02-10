/*
Package report is a utility for the report command.
*/
package report

import (
	. "github.com/jotaen/klog/src"
	"github.com/jotaen/klog/src/app/cli/lib/terminalformat"
	"github.com/jotaen/klog/src/service/period"
)

type Aggregator interface {
	NumberOfPrefixColumns() int
	DateHash(Date) period.Hash
	OnHeaderPrefix(*terminalformat.Table)
	OnRowPrefix(*terminalformat.Table, Date)
}
