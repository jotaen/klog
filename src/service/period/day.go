package period

import . "github.com/jotaen/klog/src"

type DayHash Hash

type Day struct {
	date Date
}

func NewDayFromDate(d Date) Day {
	return Day{d}
}

func (d Day) Hash() DayHash {
	hash := newBitMask()
	hash.populate(uint32(d.date.Day()), 31)
	hash.populate(uint32(d.date.Month()), 12)
	hash.populate(uint32(d.date.Year()), 10000)
	return DayHash(hash.Value())
}
