package cli

import (
	"fmt"
	"klog/app"
	. "klog/lib/tf"
	"klog/parser"
	"klog/record"
)

var Print = Command{
	Name:        "print",
	Description: "Print a file",
	Main:        print,
}

func print(ctx app.Context, args []string) int {
	if len(args) == 0 {
		fmt.Println("Please specify a file")
		return INVALID_CLI_ARGS
	}
	rs, err := ctx.Read("../" + args[0])
	if err != nil {
		fmt.Println(err)
		return EXECUTION_FAILED
	}
	h := printHooks{}
	fmt.Println(parser.SerialiseRecords(rs, h))
	return OK
}

type printHooks struct{}

func (h printHooks) PrintDate(d record.Date) string {
	return Style{Color: "222", IsUnderlined: true}.Format(d.ToString())
}
func (h printHooks) PrintShouldTotal(d record.Duration) string {
	return Style{Color: "213"}.Format(d.ToString())
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
