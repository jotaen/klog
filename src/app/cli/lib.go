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
		for _, e := range pe.Get() {
			message += fmt.Sprintf(
				Style{Color: "160"}.Format("▶︎ Syntax error in line %d:\n"),
				e.Context().LineNumber,
			)
			message += fmt.Sprintf(
				Style{Color: "247"}.Format("  %s\n"),
				string(e.Context().Value),
			)
			message += fmt.Sprintf(
				Style{Color: "160"}.Format("  %s%s\n"),
				strings.Repeat(" ", e.Position()), strings.Repeat("^", e.Length()),
			)
			message += fmt.Sprintf(
				Style{Color: "214"}.Format("  %s\n\n"),
				e.Message(),
			)
		}
		return nil, errors.New(message)
	}
	return nil, err
}

type cliPrinter struct{}

func (h cliPrinter) PrintDate(d record.Date) string {
	return Style{Color: "222", IsUnderlined: true}.Format(d.ToString())
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
