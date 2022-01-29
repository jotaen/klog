package period

import . "github.com/jotaen/klog/src"

type Week struct {
	date Date
}

type WeekHash Hash

func NewWeekFromDate(d Date) Week {
	return Week{d}
}

func (w Week) Period() Period {
	since := w.date
	until := w.date
	for {
		if since.Weekday() == 1 {
			break
		}
		since = since.PlusDays(-1)
	}
	for {
		if until.Weekday() == 7 {
			break
		}
		until = until.PlusDays(1)
	}
	return Period{since, until}
}

func (w Week) Previous() Week {
	return NewWeekFromDate(w.date.PlusDays(-7))
}

func (w Week) Hash() WeekHash {
	hash := newBitMask()
	year, week := w.date.WeekNumber()
	hash.populate(uint32(week), 53)
	hash.populate(uint32(year), 10000)
	return WeekHash(hash.Value())
}
