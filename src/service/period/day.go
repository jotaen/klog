package period

import . "github.com/jotaen/klog/src"

type DayHash Hash

type Day struct {
	Date
}

func NewDayFromDate(d Date) Day {
	return Day{d}
}

func (d Day) Hash() DayHash {
	hash := newBitMask()
	hash.populate(uint32(d.Day()), 31)
	hash.populate(uint32(d.Month()), 12)
	hash.populate(uint32(d.Year()), 10000)
	return DayHash(hash.Value())
}
