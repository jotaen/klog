package report

import (
	"github.com/jotaen/klog/lib/jotaen/terminalformat"
	. "github.com/jotaen/klog/src"
)

type Hash uint32

type Aggregator interface {
	NumberOfPrefixColumns() int
	DateHash(Date) Hash
	OnHeaderPrefix(*terminalformat.Table)
	OnRowPrefix(*terminalformat.Table, Date)
}
