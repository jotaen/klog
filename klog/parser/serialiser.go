package parser

import (
	"github.com/jotaen/klog/klog"
	"strings"
)

// Serialiser is used when the output should be modified, e.g. coloured.
type Serialiser interface {
	Date(klog.Date) string
	ShouldTotal(klog.Duration) string
	Summary(SummaryText) string
	Range(klog.Range) string
	OpenRange(klog.OpenRange) string
	Duration(klog.Duration) string
	SignedDuration(klog.Duration) string
	Time(klog.Time) string
}

type Line struct {
	Text   string
	Record klog.Record
	EntryI int
}

type Lines []Line

var canonicalLineEnding = "\n"
var canonicalIndentation = "    "

func (ls Lines) ToString() string {
	result := ""
	for _, l := range ls {
		result += l.Text + canonicalLineEnding
	}
	return result
}

// SerialiseRecords serialises records into the canonical string representation.
// (So it doesn’t and cannot restore the original formatting!)
func SerialiseRecords(s Serialiser, rs ...klog.Record) Lines {
	var lines []Line
	for i, r := range rs {
		lines = append(lines, serialiseRecord(s, r)...)
		if i < len(rs)-1 {
			lines = append(lines, Line{"", nil, -1})
		}
	}
	return lines
}

func serialiseRecord(s Serialiser, r klog.Record) []Line {
	var lines []Line
	headline := s.Date(r.Date())
	if r.ShouldTotal().InMinutes() != 0 {
		headline += " (" + s.ShouldTotal(r.ShouldTotal()) + ")"
	}
	lines = append(lines, Line{headline, r, -1})
	for _, l := range r.Summary().Lines() {
		lines = append(lines, Line{s.Summary([]string{l}), r, -1})
	}
	for entryI, e := range r.Entries() {
		entryValue := klog.Unbox[string](&e,
			func(r klog.Range) string { return s.Range(r) },
			func(d klog.Duration) string { return s.Duration(d) },
			func(o klog.OpenRange) string { return s.OpenRange(o) },
		)
		lines = append(lines, Line{canonicalIndentation + entryValue, r, entryI})
		for i, l := range e.Summary().Lines() {
			summaryText := s.Summary([]string{l})
			if i == 0 && l != "" {
				lines[len(lines)-1].Text += " " + summaryText
			} else if i >= 1 {
				lines = append(lines, Line{canonicalIndentation + canonicalIndentation + summaryText, r, entryI})
			}
		}
	}
	return lines
}

type SummaryText []string

func (s SummaryText) ToString() string {
	return strings.Join(s, canonicalLineEnding)
}
