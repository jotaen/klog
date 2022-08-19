package period

import (
	"errors"
	"github.com/jotaen/klog/klog"
	"regexp"
	"strconv"
	"strings"
)

type Week struct {
	date klog.Date
}

type WeekHash Hash

var weekPattern = regexp.MustCompile(`^\d{4}-W\d{1,2}$`)

func NewWeekFromDate(d klog.Date) Week {
	return Week{d}
}

func NewWeekFromString(yyyyWww string) (Week, error) {
	if !weekPattern.MatchString(yyyyWww) {
		return Week{}, errors.New("INVALID_WEEK_PERIOD")
	}
	parts := strings.Split(yyyyWww, "-")
	year, _ := strconv.Atoi(parts[0])
	week, _ := strconv.Atoi(strings.TrimPrefix(parts[1], "W"))
	if week < 1 {
		return Week{}, errors.New("INVALID_WEEK_PERIOD")
	}
	reference, err := func() (klog.Date, error) {
		ref, yErr := klog.NewDate(year, 7, 1)
		if yErr != nil {
			return nil, errors.New("INVALID_WEEK_PERIOD")
		}
		for ref.Weekday() != 1 {
			ref = ref.PlusDays(-1)
		}
		_, w := ref.WeekNumber()
		ref = ref.PlusDays((week - w) * 7)
		return ref, nil
	}()
	if err != nil {
		return Week{}, err
	}
	if _, refWeekNr := reference.WeekNumber(); refWeekNr != week {
		// Prevent implicit roll over.
		return Week{}, errors.New("INVALID_WEEK_PERIOD")
	}
	return Week{reference}, nil
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
	return NewPeriod(since, until)
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
