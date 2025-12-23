package kql

import "github.com/jotaen/klog/klog"

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
	ENTRY_TYPE_DURATION          = EntryType("DURATION")
	ENTRY_TYPE_POSITIVE_DURATION = EntryType("DURATION_POSITIVE")
	ENTRY_TYPE_NEGATIVE_DURATION = EntryType("DURATION_NEGATIVE")
	ENTRY_TYPE_RANGE             = EntryType("RANGE")
	ENTRY_TYPE_OPEN_RANGE        = EntryType("OPEN_RANGE")
)

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
