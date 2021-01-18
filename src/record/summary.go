package record

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

type TagSet map[Tag]bool

func NewTag(value string) Tag {
	return Tag(strings.ToLower(value))
}

func ContainsOneOfTags(tags TagSet, searchText string) bool {
	for _, m := range HashTagPattern.FindAllStringSubmatch(searchText, -1) {
		tag := NewTag(m[1])
		if tags[tag] == true {
			return true
		}
	}
	return false
}

func NewTagSet(tags ...string) TagSet {
	result := map[Tag]bool{}
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
