package klog

import (
	"errors"
	"regexp"
	"sort"
	"strings"
)

// RecordSummary contains the summary lines of the overall summary that
// appears underneath the date of a record.
type RecordSummary []string

// EntrySummary contains the summary line that appears behind the time value
// of an entry.
type EntrySummary []string

// NewRecordSummary creates a new RecordSummary from individual lines of text.
// None of the lines can start with blank characters, and none of the lines
// can be empty or blank.
func NewRecordSummary(line ...string) (RecordSummary, error) {
	for _, l := range line {
		if len(l) == 0 || regexp.MustCompile(`^[\p{Zs}\t]`).MatchString(l) {
			return nil, errors.New("MALFORMED_SUMMARY")
		}
	}
	return line, nil
}

// NewEntrySummary creates an EntrySummary from individual lines of text.
// Except for the first line, none of the lines can be empty or blank.
func NewEntrySummary(line ...string) (EntrySummary, error) {
	for i, l := range line {
		if i == 0 {
			continue
		}
		if len(l) == 0 || regexp.MustCompile(`^[\p{Zs}\t]*$`).MatchString(l) {
			return nil, errors.New("MALFORMED_SUMMARY")
		}
	}
	return line, nil
}

func (s RecordSummary) Lines() []string {
	return s
}

func (s EntrySummary) Lines() []string {
	return RecordSummary(s).Lines()
}

func (s RecordSummary) Tags() TagSet {
	tags := NewTagSet()
	for _, l := range s {
		for _, m := range HashTagPattern.FindAllStringSubmatch(l, -1) {
			tag := NewTag(m[1])
			tags[tag] = true
		}
	}
	return tags
}

func (s EntrySummary) Tags() TagSet {
	return RecordSummary(s).Tags()
}

func (s RecordSummary) Equals(summary RecordSummary) bool {
	if len(s) != len(summary) {
		return false
	}
	for i, l := range s {
		if l != summary[i] {
			return false
		}
	}
	return true
}

func (s EntrySummary) Equals(summary EntrySummary) bool {
	if len(s) == 1 && s[0] == "" && summary == nil {
		// In the case of entry summary, an empty one matches nil.
		return true
	}
	return RecordSummary(s).Equals(RecordSummary(summary))
}

var HashTagPattern = regexp.MustCompile(`#([\p{L}\d_-]+)`)

type Tag string

func (t Tag) ToString() string {
	return "#" + string(t)
}

func (ts TagSet) ToStrings() []string {
	var tags []string
	for t := range ts {
		tags = append(tags, t.ToString())
	}
	sort.Slice(tags, func(i, j int) bool {
		return tags[i] < tags[j]
	})
	return tags
}

func (ts TagSet) Contains(queryTag string) bool {
	if !strings.HasSuffix(queryTag, "...") {
		return ts[NewTag(queryTag)]
	}
	queryBaseTag := NewTag(strings.TrimSuffix(queryTag, "..."))
	for t := range ts {
		if strings.HasPrefix(t.ToString(), queryBaseTag.ToString()) {
			return true
		}
	}
	return false
}

type TagSet map[Tag]bool

func NewTag(value string) Tag {
	if value[0] == '#' {
		value = value[1:]
	}
	return Tag(strings.ToLower(value))
}

func NewTagSet(tags ...string) TagSet {
	result := make(map[Tag]bool, len(tags))
	for _, v := range tags {
		if len(v) == 0 {
			continue
		}
		tag := NewTag(v)
		result[tag] = true
	}
	return result
}

func Merge(tagSets ...TagSet) TagSet {
	result := NewTagSet()
	for _, ts := range tagSets {
		for t := range ts {
			result[t] = true
		}
	}
	return result
}
