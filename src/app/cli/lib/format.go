package lib

import (
	"errors"
	"fmt"
	. "klog"
	"klog/app"
	. "klog/lib/jotaen/tf"
	"klog/parser"
	"klog/parser/parsing"
	"klog/service"
	"strings"
)

var Styler = parser.Serialiser{
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
	Duration: func(d Duration, forceSign bool) string {
		f := Style{Color: "120"}
		if d.InMinutes() < 0 {
			f.Color = "167"
		}
		if forceSign {
			return f.Format(d.ToStringWithSign())
		}
		return f.Format(d.ToString())
	},
	Time: func(t Time) string {
		return Style{Color: "027"}.Format(t.ToString())
	},
}

var lineBreaker = lineBreakerT{
	maxLength: 60,
	newLine:   "\n",
}

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
				lineBreaker.split(e.Message(), INDENT),
			) + "\n\n"
		}
		return errors.New(message)
	case app.Error:
		message := "Error: " + e.Error() + "\n"
		message += lineBreaker.split(e.Details(), "")
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
