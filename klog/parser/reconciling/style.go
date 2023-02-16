package reconciling

import (
	"github.com/jotaen/klog/klog"
	"github.com/jotaen/klog/klog/parser/txt"
)

// style describes the general styling and formatting preferences of a record.
type style struct {
	lineEnding                          styleProp[string]
	indentation                         styleProp[string]
	dateUseDashes                       styleProp[bool]
	timeUse24HourClock                  styleProp[bool]
	rangesUseSpacesAroundDash           styleProp[bool]
	openRangeAdditionalPlaceholderChars styleProp[int]
}

func (s *style) dateFormat() klog.DateFormat {
	f := klog.DefaultDateFormat()
	f.UseDashes = s.dateUseDashes.Get()
	return f
}

func (s *style) timeFormat() klog.TimeFormat {
	f := klog.DefaultTimeFormat()
	f.Use24HourClock = s.timeUse24HourClock.Get()
	return f
}

func (s *style) openRangeFormat() klog.OpenRangeFormat {
	f := klog.DefaultOpenRangeFormat()
	f.UseSpacesAroundDash = s.rangesUseSpacesAroundDash.Get()
	f.AdditionalPlaceholderChars = s.openRangeAdditionalPlaceholderChars.Get()
	return f
}

func determine(r klog.Record, b txt.Block) *style {
	s := defaultStyle()
	s.dateUseDashes.Set(r.Date().Format().UseDashes)
	for _, e := range r.Entries() {
		klog.Unbox[any](&e, func(r klog.Range) any {
			s.timeUse24HourClock.Set(r.Start().Format().Use24HourClock)
			s.rangesUseSpacesAroundDash.Set(r.Format().UseSpacesAroundDash)
			return nil
		}, func(d klog.Duration) any {
			return nil
		}, func(o klog.OpenRange) any {
			s.timeUse24HourClock.Set(o.Start().Format().Use24HourClock)
			s.rangesUseSpacesAroundDash.Set(o.Format().UseSpacesAroundDash)
			s.openRangeAdditionalPlaceholderChars.Set(o.Format().AdditionalPlaceholderChars)
			return nil
		})
	}
	for _, l := range b.Lines() {
		if l.Indentation() != "" {
			s.indentation.Set(l.Indentation())
			break
		}
	}
	if len(b.Lines()) > 0 && b.Lines()[0].LineEnding != "" {
		s.lineEnding.Set(b.Lines()[0].LineEnding)
	}
	return s
}

type styleProp[T any] struct {
	value      T
	isExplicit bool
}

func (p *styleProp[T]) Set(value T) {
	p.value = value
	p.isExplicit = true
}

func (p *styleProp[T]) Get() T {
	return p.value
}

// defaultStyle returns the canonical style preferences as recommended
// by the file format specification.
func defaultStyle() *style {
	return &style{
		lineEnding:                          styleProp[string]{"\n", false},
		indentation:                         styleProp[string]{"    ", false},
		dateUseDashes:                       styleProp[bool]{klog.DefaultDateFormat().UseDashes, false},
		timeUse24HourClock:                  styleProp[bool]{klog.DefaultTimeFormat().Use24HourClock, false},
		rangesUseSpacesAroundDash:           styleProp[bool]{klog.DefaultRangeFormat().UseSpacesAroundDash, false},
		openRangeAdditionalPlaceholderChars: styleProp[int]{klog.DefaultOpenRangeFormat().AdditionalPlaceholderChars, false},
	}
}

type election[T comparable] struct {
	votes map[T]int
}

func newElection[T comparable]() election[T] {
	return election[T]{make(map[T]int)}
}

// vote casts a vote for the style, but only if it’s explicit.
func (e *election[T]) vote(style styleProp[T]) {
	if !style.isExplicit {
		return
	}
	e.votes[style.value] += 1
}

// tallyUp returns the style that’s most voted for.
func (e *election[T]) tallyUp(defaultValue T) T {
	max := 0
	result := defaultValue
	for value, count := range e.votes {
		if count > max {
			max = count
			result = value
		}
	}
	return result
}

// ascertain finds the prevailing style, which is either the explicit default style,
// or the tallied-up election.
func ascertain[T comparable](e *election[T], defaultStyle styleProp[T]) styleProp[T] {
	if defaultStyle.isExplicit {
		return defaultStyle
	}
	return styleProp[T]{e.tallyUp(defaultStyle.Get()), true}
}

// elect fills all unset fields of the `defaults` style with that value
// which was encountered most often in the parsed records. Fields of the
// `base` style that had been set explicitly take precedence.
func elect(base style, rs []klog.Record, bs []txt.Block) *style {
	if len(rs) != len(bs) {
		panic("ASSERTION_ERROR")
	}
	lineEndingElection := newElection[string]()
	indentationElection := newElection[string]()
	dateUseDashes := newElection[bool]()
	timeUse24HourClock := newElection[bool]()
	rangesUseSpacesAroundDash := newElection[bool]()
	openRangeAdditionalPlaceholderChars := newElection[int]()
	for i, r := range rs {
		s := determine(r, bs[i])
		lineEndingElection.vote(s.lineEnding)
		indentationElection.vote(s.indentation)
		dateUseDashes.vote(s.dateUseDashes)
		timeUse24HourClock.vote(s.timeUse24HourClock)
		rangesUseSpacesAroundDash.vote(s.rangesUseSpacesAroundDash)
		openRangeAdditionalPlaceholderChars.vote(s.openRangeAdditionalPlaceholderChars)
	}
	return &style{
		lineEnding:                          ascertain[string](&lineEndingElection, base.lineEnding),
		indentation:                         ascertain[string](&indentationElection, base.indentation),
		dateUseDashes:                       ascertain[bool](&dateUseDashes, base.dateUseDashes),
		timeUse24HourClock:                  ascertain[bool](&timeUse24HourClock, base.timeUse24HourClock),
		rangesUseSpacesAroundDash:           ascertain[bool](&rangesUseSpacesAroundDash, base.rangesUseSpacesAroundDash),
		openRangeAdditionalPlaceholderChars: ascertain[int](&openRangeAdditionalPlaceholderChars, base.openRangeAdditionalPlaceholderChars),
	}
}
