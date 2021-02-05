package cli

import (
	"errors"
	"fmt"
	"klog"
	. "klog/lib/jotaen/tf"
	"klog/parser"
	"klog/parser/engine"
	"strings"
)

func prettifyError(err error) error {
	switch e := err.(type) {
	case engine.Errors:
		message := ""
		INDENT := "    "
		for _, e := range e.Get() {
			message += fmt.Sprintf(
				Style{Background: "160", Color: "015"}.Format(" Error in line %d: "),
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
	}
	return err
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
