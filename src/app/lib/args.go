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
	File []app.FileOrBookmarkName `arg:"" optional:"" type:"string" name:"file or bookmark" help:".klg source file(s) (if empty the bookmark is used)"`
}

type OutputFileArgs struct {
	File app.FileOrBookmarkName `arg:"" optional:"" type:"string" name:"file or bookmark" help:".klg source file (if empty the bookmark is used)"`
}

type AtDateArgs struct {
	Date      Date `name:"date" short:"d" help:"The date of the record"`
	Today     bool `name:"today" help:"Use today’s date (default)"`
	Yesterday bool `name:"yesterday" help:"Use yesterday’s date"`
	Tomorrow  bool `name:"tomorrow" help:"Use tomorrow’s date"`
}

func (args *AtDateArgs) AtDate(now gotime.Time) (Date, bool) {
	if args.Date != nil {
		return args.Date, false
	}
	today := NewDateFromGo(now) // That’s effectively/implicitly `--today`
	if args.Yesterday {
		return today.PlusDays(-1), false
	} else if args.Tomorrow {
		return today.PlusDays(1), false
	}
	return today, true
}

type AtDateAndTimeArgs struct {
	AtDateArgs
	Time Time `name:"time" short:"t" help:"Specify the time (defaults to now)"`
}

func (args *AtDateAndTimeArgs) AtTime(now gotime.Time) (Time, bool, app.Error) {
	if args.Time != nil {
		return args.Time, false, nil
	}
	date, _ := args.AtDate(now)
	today := NewDateFromGo(now)
	if today.IsEqualTo(date) {
		return NewTimeFromGo(now), true, nil
	} else if today.PlusDays(-1).IsEqualTo(date) {
		shiftedTime, _ := NewTimeFromGo(now).Plus(NewDuration(24, 0))
		return shiftedTime, true, nil
	} else if today.PlusDays(1).IsEqualTo(date) {
		shiftedTime, _ := NewTimeFromGo(now).Plus(NewDuration(-24, 0))
		return shiftedTime, true, nil
	}
	return nil, false, app.NewErrorWithCode(
		app.LOGICAL_ERROR,
		"Missing time parameter",
		"Please specify a time value for dates in the past",
		nil,
	)
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
		qry.Dates = append(qry.Dates, NewDateFromGo(now))
	}
	if args.Yesterday {
		qry.Dates = append(qry.Dates, NewDateFromGo(now.AddDate(0, 0, -1)))
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
	Quiet bool `name:"quiet" help:"Output parseable data without descriptive text"`
}

type SortArgs struct {
	Sort string `name:"sort" help:"Sort output by date (asc or desc)" enum:"asc,desc,ASC,DESC," default:""`
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
