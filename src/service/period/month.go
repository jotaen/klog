package period

import (
	"errors"
	. "github.com/jotaen/klog/src"
	"regexp"
	"strconv"
	"strings"
)

type Month struct {
	year  int
	month int
}

type MonthHash Hash

var monthPattern = regexp.MustCompile(`^\d{4}-\d{2}$`)

func NewMonthFromDate(d Date) Month {
	return Month{d.Year(), d.Month()}
}

func NewMonthFromString(yyyymm string) (Month, error) {
	if !monthPattern.MatchString(yyyymm) {
		return Month{}, errors.New("INVALID_MONTH_PERIOD")
	}
	parts := strings.Split(yyyymm, "-")
	year, _ := strconv.Atoi(parts[0])
	month, _ := strconv.Atoi(parts[1])
	_, err := NewDate(year, month, 1)
	if err != nil {
		return Month{}, errors.New("INVALID_MONTH_PERIOD")
	}
	return Month{year, month}, nil
}

func (m Month) Period() Period {
	since, _ := NewDate(m.year, m.month, 1)
	until, _ := NewDate(m.year, m.month, 28)
	for {
		if until.Year() == 9999 && until.Month() == 12 && until.Day() == 31 {
			// 9999-12-31 is the last representable date, so we can’t peak forward from it.
			break
		}
		next := until.PlusDays(1)
		if next.Month() != until.Month() {
			break
		}
		until = next
	}
	return Period{
		Since: since,
		Until: until,
	}
}

func (m Month) Hash() MonthHash {
	hash := newBitMask()
	hash.populate(uint32(m.month), 12)
	hash.populate(uint32(m.year), 10000)
	return MonthHash(hash.Value())
}
