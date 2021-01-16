package record

import (
	"regexp"
	"strings"
)

func Total(r Record) Duration {
	total := NewDuration(0, 0)
	for _, e := range r.Entries() {
		switch v := e.Value().(type) {
		case Duration:
			total = total.Add(v)
			break
		case Range:
			total = total.Add(v.Duration())
			break
		}
	}
	return total
}

func Find(date Date, rs []Record) Record {
	for _, r := range rs {
		if r.Date() == date {
			return r
		}
	}
	return nil
}

var hashTagPattern = regexp.MustCompile(`#(\p{L}+)`)

func FindEntriesWithHashtags(tags map[string]bool, r Record) []Entry {
	if ContainsOneOfTags(tags, r.Summary()) {
		return r.Entries()
	}
	var matches []Entry
	for _, e := range r.Entries() {
		if ContainsOneOfTags(tags, e.SummaryAsString()) {
			matches = append(matches, e)
		}
	}
	return matches
}

func ContainsOneOfTags(tags map[string]bool, searchText string) bool {
	for _, t := range hashTagPattern.FindAllStringSubmatch(searchText, -1) {
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
