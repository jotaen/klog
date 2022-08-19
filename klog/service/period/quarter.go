package period

import (
	"errors"
	"github.com/jotaen/klog/klog"
	"regexp"
	"strconv"
	"strings"
)

type Quarter struct {
	date klog.Date
}

type QuarterHash Hash

var quarterPattern = regexp.MustCompile(`^\d{4}-Q\d$`)

func NewQuarterFromDate(d klog.Date) Quarter {
	return Quarter{d}
}

func NewQuarterFromString(yyyyQq string) (Quarter, error) {
	if !quarterPattern.MatchString(yyyyQq) {
		return Quarter{}, errors.New("INVALID_QUARTER_PERIOD")
	}
	parts := strings.Split(yyyyQq, "-")
	year, _ := strconv.Atoi(parts[0])
	quarter, _ := strconv.Atoi(strings.TrimPrefix(parts[1], "Q"))
	if quarter < 1 || quarter > 4 {
		return Quarter{}, errors.New("INVALID_QUARTER_PERIOD")
	}
	month := quarter * 3
	d, err := klog.NewDate(year, month, 1)
	if err != nil {
		return Quarter{}, errors.New("INVALID_QUARTER_PERIOD")
	}
	return Quarter{d}, nil
}

func (q Quarter) Period() Period {
	switch q.date.Quarter() {
	case 1:
		since, _ := klog.NewDate(q.date.Year(), 1, 1)
		until, _ := klog.NewDate(q.date.Year(), 3, 31)
		return NewPeriod(since, until)
	case 2:
		since, _ := klog.NewDate(q.date.Year(), 4, 1)
		until, _ := klog.NewDate(q.date.Year(), 6, 30)
		return NewPeriod(since, until)
	case 3:
		since, _ := klog.NewDate(q.date.Year(), 7, 1)
		until, _ := klog.NewDate(q.date.Year(), 9, 30)
		return NewPeriod(since, until)
	case 4:
		since, _ := klog.NewDate(q.date.Year(), 10, 1)
		until, _ := klog.NewDate(q.date.Year(), 12, 31)
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
