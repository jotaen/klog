package parser

import . "github.com/jotaen/klog/src"

// Style describes the general styling and formatting preferences of a record.
type Style struct {
	LineEnding     styleProp[string]
	Indentation    styleProp[string]
	DateFormat     styleProp[DateFormat]
	TimeFormat     styleProp[TimeFormat]
	SpacingInRange styleProp[string] // Example: `8:00 - 9:00` vs. `8:00-9:00`
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

// DefaultStyle returns the canonical style preferences as recommended
// by the file format specification.
func DefaultStyle() *Style {
	return &Style{
		LineEnding:     styleProp[string]{"\n", false},
		Indentation:    styleProp[string]{"    ", false},
		DateFormat:     styleProp[DateFormat]{DateFormat{UseDashes: true}, false},
		TimeFormat:     styleProp[TimeFormat]{TimeFormat{Use24HourClock: true}, false},
		SpacingInRange: styleProp[string]{" ", false},
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

// Elect fills all unset fields of the `defaults` style with that value
// which was encountered most often in the parsed records. Fields of the
// `base` style that had been set explicitly take precedence.
func Elect(base Style, parsedRecords []ParsedRecord) *Style {
	lineEndingElection := newElection[string]()
	indentationElection := newElection[string]()
	dateFormatElection := newElection[DateFormat]()
	timeFormatElection := newElection[TimeFormat]()
	spacingInRangeElection := newElection[string]()
	for _, r := range parsedRecords {
		lineEndingElection.vote(r.Style.LineEnding)
		indentationElection.vote(r.Style.Indentation)
		dateFormatElection.vote(r.Style.DateFormat)
		timeFormatElection.vote(r.Style.TimeFormat)
		spacingInRangeElection.vote(r.Style.SpacingInRange)
	}
	return &Style{
		ascertain[string](&lineEndingElection, base.LineEnding),
		ascertain[string](&indentationElection, base.Indentation),
		ascertain[DateFormat](&dateFormatElection, base.DateFormat),
		ascertain[TimeFormat](&timeFormatElection, base.TimeFormat),
		ascertain[string](&spacingInRangeElection, base.SpacingInRange),
	}
}
