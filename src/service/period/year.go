package period

import (
	"errors"
	. "github.com/jotaen/klog/src"
	"regexp"
	"strconv"
)

type Year struct {
	date Date
}

type YearHash Hash

var yearPattern = regexp.MustCompile(`^\d{4}$`)

func NewYearFromDate(d Date) Year {
	return Year{d}
}

func NewYearFromString(yyyy string) (Year, error) {
	if !yearPattern.MatchString(yyyy) {
		return Year{}, errors.New("INVALID_YEAR_PERIOD")
	}
	year, err := strconv.Atoi(yyyy)
	if err != nil {
		return Year{}, errors.New("INVALID_YEAR_PERIOD")
	}
	d, dErr := NewDate(year, 1, 1)
	if dErr != nil {
		return Year{}, errors.New("INVALID_YEAR_PERIOD")
	}
	return Year{d}, nil
}

func (y Year) Period() Period {
	since, _ := NewDate(y.date.Year(), 1, 1)
	until, _ := NewDate(y.date.Year(), 12, 31)
	return Period{
		Since: since,
		Until: until,
	}
}

func (y Year) Hash() YearHash {
	hash := newBitMask()
	hash.populate(uint32(y.date.Year()), 10000)
	return YearHash(hash.Value())
}
