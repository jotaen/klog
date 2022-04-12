package klog

import (
	"errors"
	"regexp"
	"sort"
	"strings"
)

var HashTagPattern = regexp.MustCompile(`#([\p{L}\d_-]+)(=(("[^"]*")|('[^']*')|([\p{L}\d_-]*)))?`)
var unquotedValuePattern = regexp.MustCompile(`^[\p{L}\d_-]+$`)

type Tag struct {
	name  string
	value string
}

func NewTagFromString(tag string) (Tag, error) {
	if !strings.HasPrefix(tag, "#") {
		tag = "#" + tag
	}
	match := HashTagPattern.FindStringSubmatch(tag)
	if match == nil {
		// The tag pattern didn’t match at all.
		return Tag{}, errors.New("INVALID_TAG")
	}
	name := match[1]
	value := func() string {
		v := match[3]
		if strings.HasPrefix(v, `"`) {
			return strings.Trim(v, `"`)
		}
		if strings.HasPrefix(v, `'`) {
			return strings.Trim(v, `'`)
		}
		return v
	}()
	if len(match[0]) != len(tag) {
		// The original tag contains more/other characters.
		return Tag{}, errors.New("INVALID_TAG")
	}
	return NewTagOrPanic(name, value), nil
}

// NewTagOrPanic constructs a new tag but will panic if the
// parameters don’t yield a valid tag.
func NewTagOrPanic(name string, value string) Tag {
	if strings.Contains(value, "\"") && strings.Contains(value, "'") {
		// A tag value can never contain both ' and " at the same time.
		panic("Invalid tag")
	}
	return Tag{strings.ToLower(name), value}
}

func (t Tag) Name() string {
	return t.name
}

func (t Tag) Value() string {
	return t.value
}

func (t Tag) Matches(queryTag Tag) bool {
	if t.Name() != queryTag.Name() {
		return false
	}
	if queryTag.Value() != "" {
		return t.Value() == queryTag.Value()
	}
	return true
}

func (t Tag) ToString() string {
	result := "#" + t.name
	if t.value != "" {
		result += "="
		quotation := ""
		if !unquotedValuePattern.MatchString(t.value) {
			if strings.Contains(t.value, `"`) {
				quotation = `'`
			} else {
				quotation = "\""
			}
		}
		result += quotation + t.value + quotation
	}
	return result
}

type TagSet map[Tag]bool

func NewEmptyTagSet() TagSet {
	return make(map[Tag]bool)
}

func (ts TagSet) Put(tag Tag) {
	ts[tag] = true
	ts[NewTagOrPanic(tag.Name(), "")] = true
}

func (ts TagSet) Contains(tag Tag) bool {
	return ts[tag]
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

func Merge(tagSets ...TagSet) TagSet {
	result := NewEmptyTagSet()
	for _, ts := range tagSets {
		for t := range ts {
			result[t] = true
		}
	}
	return result
}
