package lib

import (
	"errors"
	"fmt"
	. "klog"
	"klog/app"
	. "klog/lib/jotaen/terminalformat"
	"klog/parser"
	"klog/parser/parsing"
	"klog/service"
	"strings"
)

func NewCliSerialiser() *parser.Serialiser {
	return &parser.Serialiser{
		Date: func(d Date) string {
			return Style{Color: "015", IsUnderlined: true}.Format(d.ToString())
		},
		ShouldTotal: func(d Duration) string {
			return Style{Color: "213"}.Format(d.ToString())
		},
		Summary: func(s Summary) string {
			txt := s.ToString()
			style := Style{Color: "249"}
			hashStyle := style.ChangedBold(true).ChangedColor("251")
			txt = HashTagPattern.ReplaceAllStringFunc(txt, func(h string) string {
				return hashStyle.FormatAndRestore(h, style)
			})
			return style.Format(txt)
		},
		Range: func(r Range) string {
			return Style{Color: "117"}.Format(r.ToString())
		},
		OpenRange: func(or OpenRange) string {
			return Style{Color: "027"}.Format(or.ToString())
		},
		Duration: func(d Duration) string {
			f := Style{Color: "120"}
			if d.InMinutes() < 0 {
				f.Color = "167"
			}
			return f.Format(d.ToString())
		},
		SignedDuration: func(d Duration) string {
			f := Style{Color: "120"}
			if d.InMinutes() < 0 {
				f.Color = "167"
			}
			return f.Format(d.ToStringWithSign())
		},
		Time: func(t Time) string {
			return Style{Color: "027"}.Format(t.ToString())
		},
	}
}

var reflower = NewReflower(60, "\n")

func PrettifyError(err error, isDebug bool) error {
	switch e := err.(type) {
	case parsing.Errors:
		message := ""
		INDENT := "    "
		for _, e := range e.Get() {
			message += fmt.Sprintf(
				Style{Background: "160", Color: "015"}.Format(" ERROR in line %d: "),
				e.Context().LineNumber,
			) + "\n"
			message += fmt.Sprintf(
				Style{Color: "247"}.Format(INDENT+"%s"),
				e.Context().Text,
			) + "\n"
			message += fmt.Sprintf(
				Style{Color: "160"}.Format(INDENT+"%s%s"),
				strings.Repeat(" ", e.Position()), strings.Repeat("^", e.Length()),
			) + "\n"
			message += fmt.Sprintf(
				Style{Color: "227"}.Format("%s"),
				reflower.Reflow(e.Message(), INDENT),
			) + "\n\n"
		}
		return errors.New(message)
	case app.Error:
		message := "Error: " + e.Error() + "\n"
		message += reflower.Reflow(e.Details(), "")
		if isDebug && e.Original() != nil {
			message += "\n\nOriginal Error:\n" + e.Original().Error()
		}
		return errors.New(message)
	}
	return errors.New("Error: " + err.Error())
}

func PrettifyWarnings(ws []service.Warning) string {
	result := ""
	for _, w := range ws {
		result += Style{Background: "227", Color: "000"}.Format(" WARNING ")
		result += " "
		result += Style{Color: "227"}.Format(w.Date.ToString() + ": " + w.Message)
		result += "\n"
	}
	return result
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
