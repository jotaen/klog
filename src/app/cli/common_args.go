package cli

import (
	. "klog"
	"klog/app"
	"klog/service"
	"time"
)

type InputFilesArgs struct {
	File []string `arg optional type:"existingfile" name:"file" help:".klg source file(s) (if empty the bookmark is used)"`
}

type DiffArg struct {
	Diff bool `name:"diff" short:"d" help:"Show difference between actual and should total time"`
}

type NowArgs struct {
	Now bool `name:"now" short:"n" help:"Assume open ranges to be closed at this moment"`
}

type FilterArgs struct {
	Tags      []string `name:"tag" help:"Only records (or particular entries) that match this tag"`
	Date      []Date   `name:"date" help:"Only records at this date"`
	Today     bool     `name:"today" help:"Only records at today’s date"`
	Yesterday bool     `name:"yesterday" help:"Only records at yesterday’s date"`
	AfterEq   Date     `name:"after" help:"Only records after this date (inclusive)"`
	BeforeEq  Date     `name:"before" help:"Only records before this date (inclusive)"`
}

type WarnArgs struct {
	NoWarn bool `name:"no-warn" help:"Suppress warnings about potential mistakes"`
}

type SortArgs struct {
	Sort string `name:"sort" help:"Sort output by date (ASC or DESC)" enum:"ASC,DESC,"`
}

func (args *NowArgs) total(reference time.Time, rs ...Record) Duration {
	if args.Now {
		d, _ := service.HypotheticalTotal(reference, rs...)
		return d
	}
	return service.Total(rs...)
}

func (args *WarnArgs) printWarnings(ctx app.Context, records []Record) {
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
		filter.Dates = append(filter.Dates, NewDateFromTime(time.Now()))
	}
	if args.Yesterday {
		filter.Dates = append(filter.Dates, NewDateFromTime(time.Now().AddDate(0, 0, -1)))
	}
	return filter
}
