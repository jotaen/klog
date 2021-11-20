/*
Package report is a utility for the report command.
*/
package report

import (
	"github.com/jotaen/klog/lib/jotaen/terminalformat"
	. "github.com/jotaen/klog/src"
	"github.com/jotaen/klog/src/service"
)

type Aggregator interface {
	NumberOfPrefixColumns() int
	DateHash(Date) service.Hash
	OnHeaderPrefix(*terminalformat.Table)
	OnRowPrefix(*terminalformat.Table, Date)
}
