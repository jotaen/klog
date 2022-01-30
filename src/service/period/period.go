package period

import (
	. "github.com/jotaen/klog/src"
	"math"
)

// Period is an inclusive date range.
type Period interface {
	Since() Date
	Until() Date
}

type periodData struct {
	since Date
	until Date
}

func NewPeriod(since Date, until Date) Period {
	return &periodData{since, until}
}

func (p *periodData) Since() Date {
	return p.since
}

func (p *periodData) Until() Date {
	return p.until
}

// Hash is a super type for date-related hashes. Such a hash is
// the same when two dates fall into the same bucket, e.g. the same
// year and week for WeekHash or the same year, month and day for DayHash.
// The underlying int type doesn’t have any meaning.
type Hash uint32

type bitMask struct {
	value        uint32
	bitsConsumed uint
}

func newBitMask() bitMask {
	return bitMask{0, 0}
}

func (b *bitMask) Value() Hash {
	return Hash(b.value)
}

func (b *bitMask) populate(value uint32, maxValue uint32) {
	b.value = b.value | value<<b.bitsConsumed
	maxBits := uint(math.Ceil(math.Log2(float64(maxValue)))) + 1
	b.bitsConsumed += maxBits
	if b.bitsConsumed > 32 {
		panic("Overflow")
	}
}
