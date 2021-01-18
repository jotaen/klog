package cli

import (
	"errors"
	"fmt"
	"klog/app"
	. "klog/lib/tf"
	"klog/parser/engine"
	"klog/record"
	"strings"
)

func retrieveRecords(ctx *app.Context, file string) ([]record.Record, error) {
	rs, err := ctx.Read(file)
	if err == nil {
		return rs, nil
	}
	pe, isParserErrors := err.(engine.Errors)
	if isParserErrors {
		message := ""
		INDENT := "    "
		for _, e := range pe.Get() {
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
		return nil, errors.New(message)
	}
	return nil, err
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

type cliPrinter struct{}

func (h cliPrinter) PrintDate(d record.Date) string {
	return Style{Background: "090", Color: "015"}.Format(d.ToString())
}
func (h cliPrinter) PrintShouldTotal(d record.Duration, symbol string) string {
	return Style{Color: "213"}.Format(d.ToString()) + Style{Color: "201"}.Format(symbol)
}
func (h cliPrinter) PrintSummary(s record.Summary) string {
	txt := s.ToString()
	style := Style{Color: "249"}
	hashStyle := style.ChangedBold(true).ChangedColor("251")
	txt = record.HashTagPattern.ReplaceAllStringFunc(txt, func(h string) string {
		return hashStyle.FormatAndRestore(h, style)
	})
	return style.Format(txt)
}
func (h cliPrinter) PrintRange(r record.Range) string {
	return Style{Color: "117"}.Format(r.ToString())
}
func (h cliPrinter) PrintOpenRange(or record.OpenRange) string {
	return Style{Color: "027"}.Format(or.ToString())
}
func (h cliPrinter) PrintDuration(d record.Duration) string {
	f := Style{Color: "120"}
	if d.InMinutes() < 0 {
		f.Color = "167"
	}
	return f.Format(d.ToString())
}
