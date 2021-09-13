package report

import (
	. "github.com/jotaen/klog/src"
	"github.com/jotaen/klog/src/lib/jotaen/terminalformat"
)

type Hash uint32

type Aggregator interface {
	NumberOfPrefixColumns() int
	DateHash(Date) Hash
	OnHeaderPrefix(*terminalformat.Table)
	OnRowPrefix(*terminalformat.Table, Date)
}
