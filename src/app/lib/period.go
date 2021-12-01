package lib

import (
	"errors"
	"github.com/jotaen/klog/src"
	"regexp"
	"strconv"
	"strings"
)

var periodPattern = regexp.MustCompile(`^\d{4}(-\d{2})?$`)

// Period is a representation of a period of time that is representable
// by the patterns YYYY-MM or YYYY.
type Period struct {
	Since klog.Date
	Until klog.Date
}

// NewPeriodFromString turns a string into a Period objects. The string
// must be formatted according to the patterns YYYY-MM or YYYY.
func NewPeriodFromString(yyyymm string) (Period, error) {
	if yyyymm == "" || !periodPattern.MatchString(yyyymm) {
		return Period{}, errors.New("Please provide a valid period")
	}
	parts := strings.Split(yyyymm, "-")
	year, _ := strconv.Atoi(parts[0])
	monthStart := 1
	monthEnd := 12
	if len(parts) == 2 {
		monthStart, _ = strconv.Atoi(parts[1])
		monthEnd = monthStart
	}
	start, _ := klog.NewDate(year, monthStart, 1)
	end, _ := klog.NewDate(year, monthEnd, 28)
	for {
		next := end.PlusDays(1)
		if next.Month() != end.Month() {
			break
		}
		end = next
	}
	return Period{
		Since: start,
		Until: end,
	}, nil
}
