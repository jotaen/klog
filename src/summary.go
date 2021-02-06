package klog

import (
	"regexp"
	"strings"
)

type Summary string

func (s Summary) ToString() string {
	return string(s)
}

var HashTagPattern = regexp.MustCompile(`#([\p{L}\d_]+)`)

type Tag string

func (t Tag) ToString() string {
	return "#" + string(t)
}

type TagSet map[Tag]bool

func NewTag(value string) Tag {
	return Tag(strings.ToLower(value))
}

func (s Summary) MatchTags(tags TagSet) TagSet {
	matches := NewTagSet()
	allTags := s.Tags()
	for t := range tags {
		if allTags[t] {
			matches[t] = true
		}
	}
	return matches
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
		if v[0] == '#' {
			v = v[1:]
		}
		tag := NewTag(v)
		result[tag] = true
	}
	return result
}
