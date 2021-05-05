package lib

import "strings"

func Pad(length int) string {
	if length < 0 {
		return ""
	}
	return strings.Repeat(" ", length)
}

func PrettyMonth(m int) string {
	switch m {
	case 1:
		return "January"
	case 2:
		return "February"
	case 3:
		return "March"
	case 4:
		return "April"
	case 5:
		return "May"
	case 6:
		return "June"
	case 7:
		return "July"
	case 8:
		return "August"
	case 9:
		return "September"
	case 10:
		return "October"
	case 11:
		return "November"
	case 12:
		return "December"
	}
	panic("Illegal month") // this can/should never happen
}

func PrettyDay(d int) string {
	switch d {
	case 1:
		return "Monday"
	case 2:
		return "Tuesday"
	case 3:
		return "Wednesday"
	case 4:
		return "Thursday"
	case 5:
		return "Friday"
	case 6:
		return "Saturday"
	case 7:
		return "Sunday"
	}
	panic("Illegal weekday") // this can/should never happen
}

type lineBreakerT struct {
	maxLength int
	newLine   string
}

func (b lineBreakerT) split(text string, linePrefix string) string {
	SPACE := " "
	words := strings.Split(text, SPACE)
	lines := []string{""}
	for i, word := range words {
		nr := len(lines) - 1
		isLastWordOfText := i == len(words)-1
		if !isLastWordOfText && len(lines[nr])+len(words[i+1]) > b.maxLength {
			lines = append(lines, "")
			nr = len(lines) - 1
		}
		if lines[nr] == "" {
			lines[nr] += linePrefix
		} else {
			lines[nr] += SPACE
		}
		lines[nr] += word
	}
	return strings.Join(lines, b.newLine)
}
