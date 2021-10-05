package lib

import (
	. "github.com/jotaen/klog/src"
	"github.com/jotaen/klog/src/app"
	"github.com/jotaen/klog/src/parser"
	"github.com/jotaen/klog/src/service"
	"os"
	"strings"
	gotime "time"
)

type InputFilesArgs struct {
	File []app.FileOrBookmarkName `arg optional type:"fileOrBookmarkName" name:"file or bookmark" help:".klg source file(s) (if empty the bookmark is used)"`
}

type OutputFileArgs struct {
	File app.FileOrBookmarkName `arg optional type:"fileOrBookmarkName" name:"file or bookmark" help:".klg source file (if empty the bookmark is used)"`
}

type AtDateArgs struct {
	Today     bool `name:"today" help:"Use today’s date (default)"`
	Yesterday bool `name:"yesterday" help:"Use yesterday’s date"`
	Date      Date `name:"date" short:"d" help:"The date of the record"`
}

func (args *AtDateArgs) AtDate(now gotime.Time) Date {
	if args.Date != nil {
		return args.Date
	}
	today := NewDateFromTime(now)
	if args.Yesterday {
		return today.PlusDays(-1)
	}
	return today
}

type AtTimeArgs struct {
	Time Time `name:"time" short:"t" help:"Specify the time (defaults to now)"`
}

func (args *AtTimeArgs) AtTime(now gotime.Time) Time {
	if args.Time != nil {
		return args.Time
	}
	return NewTimeFromTime(now)
}

type DiffArgs struct {
	Diff bool `name:"diff" short:"d" help:"Show difference between actual and should-total time"`
}

type NowArgs struct {
	Now bool `name:"now" short:"n" help:"Assume open ranges to be closed at this moment"`
}

func (args *NowArgs) Total(reference gotime.Time, rs ...Record) Duration {
	if args.Now {
		d, _ := service.HypotheticalTotal(reference, rs...)
		return d
	}
	return service.Total(rs...)
}

type FilterArgs struct {
	Tags      []string `name:"tag" group:"Filter" help:"Only records (or particular entries) that match this tag"`
	Date      []Date   `name:"date" group:"Filter" help:"Only records at this date"`
	Today     bool     `name:"today" group:"Filter" help:"Only records at today’s date"`
	Yesterday bool     `name:"yesterday" group:"Filter" help:"Only records at yesterday’s date"`
	Since     Date     `name:"since" group:"Filter" help:"Only records since this date (inclusive)"`
	Until     Date     `name:"until" group:"Filter" help:"Only records until this date (inclusive)"`
	After     Date     `name:"after" group:"Filter" help:"Only records after this date (exclusive)"`
	Before    Date     `name:"before" group:"Filter" help:"Only records before this date (exclusive)"`
	Period    Period   `name:"period" group:"Filter" help:"Only records in this period (YYYY-MM or YYYY)"`
}

func (args *FilterArgs) ApplyFilter(now gotime.Time, rs []Record) []Record {
	qry := service.FilterQry{
		BeforeOrEqual: args.Until,
		AfterOrEqual:  args.Since,
		Tags:          args.Tags,
		Dates:         args.Date,
	}
	if args.Period.Since != nil {
		qry.BeforeOrEqual = args.Period.Until
		qry.AfterOrEqual = args.Period.Since
	}
	if args.After != nil {
		qry.AfterOrEqual = args.After.PlusDays(1)
	}
	if args.Before != nil {
		qry.BeforeOrEqual = args.Before.PlusDays(-1)
	}
	if args.Today {
		qry.Dates = append(qry.Dates, NewDateFromTime(now))
	}
	if args.Yesterday {
		qry.Dates = append(qry.Dates, NewDateFromTime(now.AddDate(0, 0, -1)))
	}
	return service.Filter(rs, qry)
}

type WarnArgs struct {
	NoWarn bool `name:"no-warn" help:"Suppress warnings about potential mistakes"`
}

func (args *WarnArgs) ToString(now gotime.Time, records []Record) string {
	if args.NoWarn {
		return ""
	}
	ws := service.SanityCheck(now, records)
	return PrettifyWarnings(ws)
}

type NoStyleArgs struct {
	NoStyle bool `name:"no-style" help:"Do not style or color the values"`
}

func (args *NoStyleArgs) Apply(ctx *app.Context) {
	if args.NoStyle || os.Getenv("NO_COLOR") != "" {
		(*ctx).SetSerialiser(&parser.PlainSerialiser)
	}
}

type QuietArgs struct {
	Quiet bool `name:"quiet" help:"Output raw data without any labels"`
}

type SortArgs struct {
	Sort string `name:"sort" help:"Sort output by date (ASC or DESC)" enum:"ASC,DESC,asc,desc,"`
}

func (args *SortArgs) ApplySort(rs []Record) []Record {
	if args.Sort == "" {
		return rs
	}
	startWithOldest := false
	if strings.ToLower(args.Sort) == "asc" {
		startWithOldest = true
	}
	return service.Sort(rs, startWithOldest)
}
