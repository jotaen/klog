package filter

import (
	"fmt"
	"strings"

	"github.com/jotaen/klog/klog"
)

// Predicate is the generic base type for all predicates. The caller is responsible for
// selecting the applicable match function (it’s meant to be either/or).
type Predicate interface {
	// Matches returns true if the record’s entry satisfies the predicate.
	Matches(klog.Record, klog.Entry) bool
	// MatchesEmptyRecord returns true if an empty record (without any entries)
	// satisfies the predicate.
	MatchesEmptyRecord(klog.Record) bool
}

type IsInDateRange struct {
	From klog.Date // May be nil to denote open range.
	To   klog.Date // May be nil to denote open range.
}

func (i IsInDateRange) Matches(r klog.Record, e klog.Entry) bool {
	return i.MatchesEmptyRecord(r)
}

func (i IsInDateRange) MatchesEmptyRecord(r klog.Record) bool {
	isAfter := func() bool {
		if i.From == nil {
			return true
		}
		return r.Date().IsAfterOrEqual(i.From)
	}()
	isBefore := func() bool {
		if i.To == nil {
			return true
		}
		return i.To.IsAfterOrEqual(r.Date())
	}()
	return isAfter && isBefore
}

type HasTag struct {
	Tag klog.Tag
}

func (h HasTag) Matches(r klog.Record, e klog.Entry) bool {
	return h.MatchesEmptyRecord(r) || e.Summary().Tags().Contains(h.Tag)
}

func (h HasTag) MatchesEmptyRecord(r klog.Record) bool {
	return r.Summary().Tags().Contains(h.Tag)
}

type And struct {
	Predicates []Predicate
}

func (a And) Matches(r klog.Record, e klog.Entry) bool {
	for _, p := range a.Predicates {
		if !p.Matches(r, e) {
			return false
		}
	}
	return true
}

func (a And) MatchesEmptyRecord(r klog.Record) bool {
	for _, p := range a.Predicates {
		if !p.MatchesEmptyRecord(r) {
			return false
		}
	}
	return true
}

type Or struct {
	Predicates []Predicate
}

func (o Or) Matches(r klog.Record, e klog.Entry) bool {
	for _, p := range o.Predicates {
		if p.Matches(r, e) {
			return true
		}
	}
	return false
}

func (o Or) MatchesEmptyRecord(r klog.Record) bool {
	for _, p := range o.Predicates {
		if p.MatchesEmptyRecord(r) {
			return true
		}
	}
	return false
}

type Not struct {
	Predicate Predicate
}

func (n Not) Matches(r klog.Record, e klog.Entry) bool {
	return !n.Predicate.Matches(r, e)
}

func (n Not) MatchesEmptyRecord(r klog.Record) bool {
	return !n.Predicate.MatchesEmptyRecord(r)
}

type EntryType string

const (
	ENTRY_TYPE_DURATION          = EntryType("duration")
	ENTRY_TYPE_DURATION_POSITIVE = EntryType("duration-positive")
	ENTRY_TYPE_DURATION_NEGATIVE = EntryType("duration-negative")
	ENTRY_TYPE_RANGE             = EntryType("range")
	ENTRY_TYPE_OPEN_RANGE        = EntryType("open-range")
)

func NewEntryTypeFromString(val string) (EntryType, error) {
	for _, t := range []EntryType{
		ENTRY_TYPE_DURATION,
		ENTRY_TYPE_DURATION_POSITIVE,
		ENTRY_TYPE_DURATION_NEGATIVE,
		ENTRY_TYPE_RANGE,
		ENTRY_TYPE_OPEN_RANGE,
	} {
		if strings.ToLower(strings.ReplaceAll(val, "_", "-")) == string(t) {
			return t, nil
		}
	}
	return EntryType(""), fmt.Errorf("%s is not a valid entry type", val)
}

type IsEntryType struct {
	Type EntryType
}

func (t IsEntryType) Matches(r klog.Record, e klog.Entry) bool {
	return klog.Unbox[bool](&e, func(r klog.Range) bool {
		return t.Type == ENTRY_TYPE_RANGE
	}, func(duration klog.Duration) bool {
		if t.Type == ENTRY_TYPE_DURATION {
			return true
		}
		if t.Type == ENTRY_TYPE_DURATION_POSITIVE && e.Duration().InMinutes() >= 0 {
			return true
		}
		if t.Type == ENTRY_TYPE_DURATION_NEGATIVE && e.Duration().InMinutes() < 0 {
			return true
		}
		return false
	}, func(openRange klog.OpenRange) bool {
		return t.Type == ENTRY_TYPE_OPEN_RANGE
	})
}

func (t IsEntryType) MatchesEmptyRecord(r klog.Record) bool {
	return false
}
