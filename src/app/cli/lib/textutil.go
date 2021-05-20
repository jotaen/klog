package lib

import (
	"strings"
)

// Deprecated
func Pad(length int) string {
	if length < 0 {
		return ""
	}
	return strings.Repeat(" ", length)
}
