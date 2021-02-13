package cli

import (
	"errors"
	"fmt"
	"klog"
	"klog/app"
	. "klog/lib/jotaen/tf"
	"klog/parser"
	"klog/parser/engine"
	"klog/service"
	"strings"
)

var styler = parser.Serialiser{
	Date: func(d klog.Date) string {
		return Style{Color: "015", IsUnderlined: true}.Format(d.ToString())
	},
	ShouldTotal: func(d klog.Duration) string {
		return Style{Color: "213"}.Format(d.ToString())
	},
	Summary: func(s klog.Summary) string {
		txt := s.ToString()
		style := Style{Color: "249"}
		hashStyle := style.ChangedBold(true).ChangedColor("251")
		txt = klog.HashTagPattern.ReplaceAllStringFunc(txt, func(h string) string {
			return hashStyle.FormatAndRestore(h, style)
		})
		return style.Format(txt)
	},
	Range: func(r klog.Range) string {
		return Style{Color: "117"}.Format(r.ToString())
	},
	OpenRange: func(or klog.OpenRange) string {
		return Style{Color: "027"}.Format(or.ToString())
	},
	Duration: func(d klog.Duration, forceSign bool) string {
		f := Style{Color: "120"}
		if d.InMinutes() < 0 {
			f.Color = "167"
		}
		if forceSign {
			return f.Format(d.ToStringWithSign())
		}
		return f.Format(d.ToString())
	},
}

func pad(length int) string {
	if length < 0 {
		return ""
	}
	return strings.Repeat(" ", length)
}

func prettyMonth(m int) string {
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

func prettyDay(d int) string {
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

func breakLines(text string, maxLength int) []string {
	SPACE := " "
	words := strings.Split(text, SPACE)
	lines := []string{""}
	for i, w := range words {
		lastLine := lines[len(lines)-1]
		isLastWord := i == len(words)-1
		if !isLastWord && len(lastLine)+len(words[i+1]) > maxLength {
			lines = append(lines, "")
		}
		lines[len(lines)-1] += w + SPACE
	}
	return lines
}

func prettifyError(err error) error {
	switch e := err.(type) {
	case engine.Errors:
		message := ""
		INDENT := "    "
		for _, e := range e.Get() {
			message += fmt.Sprintf(
				Style{Background: "160", Color: "015"}.Format(" ERROR in line %d: "),
				e.Context().LineNumber,
			) + "\n"
			message += fmt.Sprintf(
				Style{Color: "247"}.Format(INDENT+"%s"),
				string(e.Context().Value),
			) + "\n"
			message += fmt.Sprintf(
				Style{Color: "160"}.Format(INDENT+"%s%s"),
				strings.Repeat(" ", e.Position()), strings.Repeat("^", e.Length()),
			) + "\n"
			message += fmt.Sprintf(
				Style{Color: "227"}.Format(INDENT+"%s"),
				strings.Join(breakLines(e.Message(), 60), "\n"+INDENT),
			) + "\n\n"
		}
		return errors.New(message)
	case app.Error:
		return errors.New("Error: " + e.Error() + "\n" + e.Help())
	}
	return errors.New("Error: " + err.Error())
}

func prettifyWarnings(ws []service.Warning) string {
	result := ""
	for _, w := range ws {
		result += Style{Background: "227", Color: "000"}.Format(" WARNING ")
		result += " "
		result += Style{Color: "227"}.Format(w.Date.ToString() + ": " + w.Message)
		result += "\n"
	}
	return result
}
