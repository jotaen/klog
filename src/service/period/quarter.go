package period

import . "github.com/jotaen/klog/src"

type Quarter struct {
	date Date
}

type QuarterHash Hash

func NewQuarterFromDate(d Date) Quarter {
	return Quarter{d}
}

func (q Quarter) Hash() QuarterHash {
	hash := newBitMask()
	hash.populate(uint32(q.date.Quarter()), 4)
	hash.populate(uint32(q.date.Year()), 10000)
	return QuarterHash(hash.Value())
}
