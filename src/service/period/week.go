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
	year, week := w.date.WeekNumber()
	hash.populate(uint32(week), 53)
	hash.populate(uint32(year), 10000)
	return WeekHash(hash.Value())
}
