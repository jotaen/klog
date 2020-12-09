package datetime

import (
	"errors"
	"fmt"
	"regexp"
	"strconv"
)

type Duration int // in minutes

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

func (d Duration) ToString() string {
	hours := abs(int((int(d) / 60)))
	minutes := abs(int(d) % 60)
	result := ""
	if int(d) < 0 {
		result += "-"
	}
	if hours > 0 {
		result += fmt.Sprintf("%dh", hours)
	}
	if hours > 0 && minutes > 0 {
		result += " "
	}
	if minutes > 0 {
		result += fmt.Sprintf("%dm", minutes)
	}
	return result
}

var pattern = regexp.MustCompile(`^ *(-)?((\d+)h)? *((\d+)m)? *$`)

func CreateDurationFromString(hhmm string) (Duration, error) {
	match := pattern.FindStringSubmatch(hhmm)
	if match == nil {
		return 0, errors.New("INVALID_DURATION")
	}
	sign := 1
	if match[1] == "-" {
		sign = -1
	}
	hours, _ := strconv.Atoi(match[3])
	minutes, _ := strconv.Atoi(match[5])
	if minutes > 60 {
		return 0, errors.New("INVALID_DURATION")
	}
	return Duration(sign * (hours*60 + minutes)), nil
}
