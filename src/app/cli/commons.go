package cli

import (
	"errors"
	"fmt"
	"klog/app"
	"klog/lib/tf"
	"klog/parser/engine"
	"klog/record"
	"klog/service"
	"strings"
)

type FileArgs struct {
	File string `arg optional name:"file" help:"File to read from"`
}

type FilterArgs struct {
	Tags     []string    `short:"t" name:"tag" help:"Only records that contain this tag"`
	Date     record.Date `short:"d" name:"date" help:"Only records at this date"`
	AfterEq  record.Date `short:"a" name:"after" help:"Only records at or after this date"`
	BeforeEq record.Date `short:"b" name:"before" help:"Only records at or before this date"`
}

func (args *FilterArgs) toFilter() service.Filter {
	filter := service.Filter{
		BeforeEq: args.BeforeEq,
		AfterEq:  args.AfterEq,
		Tags:     args.Tags,
	}
	if args.Date != nil {
		filter.BeforeEq = args.Date
		filter.AfterEq = args.Date
	}
	return filter
}

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
				tf.Style{Color: "160"}.Format("▶︎ Syntax error in line %d:\n"),
				e.Context().LineNumber,
			)
			message += fmt.Sprintf(
				tf.Style{Color: "247"}.Format("  %s\n"),
				string(e.Context().Value),
			)
			message += fmt.Sprintf(
				tf.Style{Color: "160"}.Format("  %s%s\n"),
				strings.Repeat(" ", e.Position()), strings.Repeat("^", e.Length()),
			)
			message += fmt.Sprintf(
				tf.Style{Color: "214"}.Format("  %s\n\n"),
				e.Message(),
			)
		}
		return nil, errors.New(message)
	}
	return nil, err
}
