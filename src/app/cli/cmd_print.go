package cli

import (
	"errors"
	"fmt"
	"klog/app"
	. "klog/lib/tf"
	"klog/parser"
	"klog/record"
	"klog/service"
)

type Print struct {
	FilterArgs
	FileArgs
}

func (args *Print) Run(ctx *app.Context) error {
	rs, err := ctx.Read(args.File)
	if err != nil {
		return errors.New("EXECUTION_FAILED")
	}
	rs, _ = service.FindFilter(rs, args.FilterArgs.ToFilter())
	h := printHooks{}
	fmt.Println(parser.SerialiseRecords(rs, h))
	return nil
}

type printHooks struct{}

func (h printHooks) PrintDate(d record.Date) string {
	return Style{Color: "222", IsUnderlined: true}.Format(d.ToString())
}
func (h printHooks) PrintShouldTotal(d record.Duration, symbol string) string {
	return Style{Color: "213"}.Format(d.ToString()) + Style{Color: "201"}.Format(symbol)
}
func (h printHooks) PrintSummary(s record.Summary) string {
	txt := s.ToString()
	style := Style{Color: "249"}
	hashStyle := style.ChangedBold(true).ChangedColor("251")
	txt = record.HashTagPattern.ReplaceAllStringFunc(txt, func(h string) string {
		return hashStyle.FormatAndRestore(h, style)
	})
	return style.Format(txt)
}
func (h printHooks) PrintRange(r record.Range) string {
	return Style{Color: "117"}.Format(r.ToString())
}
func (h printHooks) PrintOpenRange(or record.OpenRange) string {
	return Style{Color: "027"}.Format(or.ToString())
}
func (h printHooks) PrintDuration(d record.Duration) string {
	f := Style{Color: "120"}
	if d.InMinutes() < 0 {
		f.Color = "167"
	}
	return f.Format(d.ToString())
}
