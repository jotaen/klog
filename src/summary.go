package klog

import (
	"regexp"
	"sort"
	"strings"
)

// Summary is arbitrary text that can be associated with a Record or an Entry.
type Summary string

func (s Summary) ToString() string {
	return string(s)
}

var HashTagPattern = regexp.MustCompile(`#([\p{L}\d_]+)`)

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

func (s Summary) Tags() TagSet {
	tags := NewTagSet()
	for _, m := range HashTagPattern.FindAllStringSubmatch(string(s), -1) {
		tag := NewTag(m[1])
		tags[tag] = true
	}
	return tags
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
