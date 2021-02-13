package cli

import (
	"klog"
	"klog/app"
	"klog/service"
	"time"
)

type InputFilesArgs struct {
	File []string `arg optional type:"existingfile" name:"file" help:".klg source file(s) (if empty the bookmark is used)"`
}

type DiffArg struct {
	Diff bool `name:"diff" help:"Show difference between actual and should total time"`
}

type FilterArgs struct {
	Tags      []string    `name:"tag" help:"Only records (or particular entries) that match this tag"`
	Date      []klog.Date `name:"date" help:"Only records at this date"`
	Today     bool        `name:"today" help:"Only records at today’s date"`
	Yesterday bool        `name:"yesterday" help:"Only records at yesterday’s date"`
	AfterEq   klog.Date   `name:"after" help:"Only records after this date (inclusive)"`
	BeforeEq  klog.Date   `name:"before" help:"Only records before this date (inclusive)"`
}

type WarnArgs struct {
	NoWarn bool `name:"nowarn" help:"Suppress warnings about potential mistakes"`
}

type SortArgs struct {
	Sort string `name:"sort" help:"Sort output by date (ASC or DESC)" enum:"ASC,DESC,"`
}

func (args *WarnArgs) printWarnings(ctx app.Context, records []klog.Record) {
	if args.NoWarn {
		return
	}
	ws := service.SanityCheck(time.Now(), records)
	ctx.Print(prettifyWarnings(ws))
}

func (args *FilterArgs) toFilter() service.Opts {
	filter := service.Opts{
		BeforeEq: args.BeforeEq,
		AfterEq:  args.AfterEq,
		Tags:     args.Tags,
		Dates:    args.Date,
	}
	if args.Today {
		filter.Dates = append(filter.Dates, klog.NewDateFromTime(time.Now()))
	}
	if args.Yesterday {
		filter.Dates = append(filter.Dates, klog.NewDateFromTime(time.Now().AddDate(0, 0, -1)))
	}
	return filter
}
