package report

import (
	. "klog"
	"klog/lib/jotaen/terminalformat"
)

type Hash uint32

type Aggregator interface {
	NumberOfPrefixColumns() int
	DateHash(Date) Hash
	OnHeaderPrefix(*terminalformat.Table)
	OnRowPrefix(*terminalformat.Table, Date)
}
