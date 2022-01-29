package period

import . "github.com/jotaen/klog/src"

type Week struct {
	date Date
}

type WeekHash Hash

func NewWeekFromDate(d Date) Week {
	return Week{d}
}

//func (w Week) Previous() Week {
//
//}

func (w Week) Hash() WeekHash {
	hash := newBitMask()
	hash.populate(uint32(w.date.WeekNumber()), 53)
	hash.populate(uint32(w.date.Year()), 10000)
	return WeekHash(hash.Value())
}
