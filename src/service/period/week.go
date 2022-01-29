package period

import . "github.com/jotaen/klog/src"

type Week struct {
	year       int
	weekNumber int
}

type WeekHash Hash

func NewWeekFromDate(d Date) Week {
	return Week{year: d.Year(), weekNumber: d.WeekNumber()}
}

func (w Week) Hash() WeekHash {
	hash := newBitMask()
	hash.populate(uint32(w.weekNumber), 53)
	hash.populate(uint32(w.year), 10000)
	return WeekHash(hash.Value())
}
