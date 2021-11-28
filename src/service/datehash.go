package service

import (
	. "github.com/jotaen/klog/src"
	"math"
)

// Hash is a super type for date-related hashes. Such a hash is
// the same when two dates fall into the same bucket, e.g. the same
// year and week for WeekHash or the same year, month and day for DayHash.
// The underlying int type doesn’t have any meaning.
type Hash uint32

type DayHash Hash
type WeekHash Hash
type MonthHash Hash
type QuarterHash Hash
type YearHash Hash

func NewDayHash(d Date) DayHash {
	hash := newBitMask()
	hash.populate(uint32(d.Day()), 31)
	hash.populate(uint32(d.Month()), 12)
	hash.populate(uint32(d.Year()), 10000)
	return DayHash(hash.Value())
}

func NewWeekHash(d Date) WeekHash {
	hash := newBitMask()
	hash.populate(uint32(d.WeekNumber()), 53)
	hash.populate(uint32(d.Year()), 10000)
	return WeekHash(hash.Value())
}

func NewMonthHash(d Date) MonthHash {
	hash := newBitMask()
	hash.populate(uint32(d.Month()), 12)
	hash.populate(uint32(d.Year()), 10000)
	return MonthHash(hash.Value())
}

func NewQuarterHash(d Date) QuarterHash {
	hash := newBitMask()
	hash.populate(uint32(d.Quarter()), 4)
	hash.populate(uint32(d.Year()), 10000)
	return QuarterHash(hash.Value())
}

func NewYearHash(d Date) YearHash {
	hash := newBitMask()
	hash.populate(uint32(d.Year()), 10000)
	return YearHash(hash.Value())
}

type bitMask struct {
	value        uint32
	bitsConsumed uint
}

func newBitMask() bitMask {
	return bitMask{0, 0}
}

func (b *bitMask) populate(value uint32, maxValue uint32) {
	b.value = b.value | value<<b.bitsConsumed
	maxBits := uint(math.Ceil(math.Log2(float64(maxValue)))) + 1
	b.bitsConsumed += maxBits
	if b.bitsConsumed > 32 {
		panic("Overflow")
	}
}

func (b *bitMask) Value() Hash {
	return Hash(b.value)
}
