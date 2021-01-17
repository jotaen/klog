package record

import (
	"regexp"
	"strings"
)

type Summary string

func (s Summary) ToString() string {
	return string(s)
}

var HashTagPattern = regexp.MustCompile(`#(\p{L}+)`)

func ContainsOneOfTags(tags map[string]bool, searchText string) bool {
	for _, t := range HashTagPattern.FindAllStringSubmatch(searchText, -1) {
		if tags[strings.ToLower(t[1])] == true {
			return true
		}
	}
	return false
}

func TagList(tags ...string) map[string]bool {
	result := map[string]bool{}
	for _, t := range tags {
		result[strings.ToLower(t)] = true
	}
	return result
}
