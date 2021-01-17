package record

import "regexp"

type Summary string

func (s Summary) ToString() string {
	return string(s)
}

var HashTagPattern = regexp.MustCompile(`#(\p{L}+)`)
