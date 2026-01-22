package filter

import (
	"errors"
	"strings"

	"github.com/jotaen/klog/klog"
)

type queriedEntry struct {
	parent klog.Record
	entry  klog.Entry
}

type Predicate interface {
	Matches(queriedEntry) bool
}

type IsInDateRange struct {
	From klog.Date
	To   klog.Date
}

func (i IsInDateRange) Matches(e queriedEntry) bool {
	isAfter := func() bool {
		if i.From == nil {
			return true
		}
		return e.parent.Date().IsAfterOrEqual(i.From)
	}()
	isBefore := func() bool {
		if i.To == nil {
			return true
		}
		return i.To.IsAfterOrEqual(e.parent.Date())
	}()
	return isAfter && isBefore
}

type HasTag struct {
	Tag klog.Tag
}

func (h HasTag) Matches(e queriedEntry) bool {
	return e.parent.Summary().Tags().Contains(h.Tag) || e.entry.Summary().Tags().Contains(h.Tag)
}

type And struct {
	Predicates []Predicate
}

func (a And) Matches(e queriedEntry) bool {
	for _, p := range a.Predicates {
		if !p.Matches(e) {
			return false
		}
	}
	return true
}

type Or struct {
	Predicates []Predicate
}

func (o Or) Matches(e queriedEntry) bool {
	for _, p := range o.Predicates {
		if p.Matches(e) {
			return true
		}
	}
	return false
}

type Not struct {
	Predicate Predicate
}

func (n Not) Matches(e queriedEntry) bool {
	return !n.Predicate.Matches(e)
}

type EntryType string

const (
	ENTRY_TYPE_DURATION          = EntryType("duration")
	ENTRY_TYPE_POSITIVE_DURATION = EntryType("duration-positive")
	ENTRY_TYPE_NEGATIVE_DURATION = EntryType("duration-negative")
	ENTRY_TYPE_RANGE             = EntryType("range")
	ENTRY_TYPE_OPEN_RANGE        = EntryType("open-range")
)

func NewEntryTypeFromString(val string) (EntryType, error) {
	for _, t := range []EntryType{
		ENTRY_TYPE_DURATION,
		ENTRY_TYPE_POSITIVE_DURATION,
		ENTRY_TYPE_NEGATIVE_DURATION,
		ENTRY_TYPE_RANGE,
		ENTRY_TYPE_OPEN_RANGE,
	} {
		if strings.ToLower(strings.ReplaceAll(val, "_", "-")) == string(t) {
			return t, nil
		}
	}
	return EntryType(""), errors.New("Illegal entry type")
}

type IsEntryType struct {
	Type EntryType
}

func (t IsEntryType) Matches(e queriedEntry) bool {
	return klog.Unbox[bool](&e.entry, func(r klog.Range) bool {
		return t.Type == ENTRY_TYPE_RANGE
	}, func(duration klog.Duration) bool {
		if t.Type == ENTRY_TYPE_DURATION {
			return true
		}
		if t.Type == ENTRY_TYPE_POSITIVE_DURATION && e.entry.Duration().InMinutes() >= 0 {
			return true
		}
		if t.Type == ENTRY_TYPE_NEGATIVE_DURATION && e.entry.Duration().InMinutes() < 0 {
			return true
		}
		return false
	}, func(openRange klog.OpenRange) bool {
		return t.Type == ENTRY_TYPE_OPEN_RANGE
	})
}
