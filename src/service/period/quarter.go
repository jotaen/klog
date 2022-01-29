package period

import . "github.com/jotaen/klog/src"

type Quarter struct {
	year    int
	quarter int
}

type QuarterHash Hash

func NewQuarterFromDate(d Date) Quarter {
	return Quarter{d.Year(), d.Quarter()}
}

func (q Quarter) Hash() QuarterHash {
	hash := newBitMask()
	hash.populate(uint32(q.quarter), 4)
	hash.populate(uint32(q.year), 10000)
	return QuarterHash(hash.Value())
}
