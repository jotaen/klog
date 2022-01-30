package period

import . "github.com/jotaen/klog/src"

type Quarter struct {
	date Date
}

type QuarterHash Hash

func NewQuarterFromDate(d Date) Quarter {
	return Quarter{d}
}

func (q Quarter) Period() Period {
	switch q.date.Quarter() {
	case 1:
		since, _ := NewDate(q.date.Year(), 1, 1)
		until, _ := NewDate(q.date.Year(), 3, 31)
		return NewPeriod(since, until)
	case 2:
		since, _ := NewDate(q.date.Year(), 4, 1)
		until, _ := NewDate(q.date.Year(), 6, 30)
		return NewPeriod(since, until)
	case 3:
		since, _ := NewDate(q.date.Year(), 7, 1)
		until, _ := NewDate(q.date.Year(), 9, 30)
		return NewPeriod(since, until)
	case 4:
		since, _ := NewDate(q.date.Year(), 10, 1)
		until, _ := NewDate(q.date.Year(), 12, 31)
		return NewPeriod(since, until)
	}
	// This can/should never happen
	panic("Invalid quarter")
}

func (q Quarter) Previous() Quarter {
	result := q.date
	for {
		result = result.PlusDays(-80)
		if result.Quarter() != q.date.Quarter() {
			return Quarter{result}
		}
	}
}

func (q Quarter) Hash() QuarterHash {
	hash := newBitMask()
	hash.populate(uint32(q.date.Quarter()), 4)
	hash.populate(uint32(q.date.Year()), 10000)
	return QuarterHash(hash.Value())
}
