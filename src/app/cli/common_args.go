package cli

import (
	. "klog"
	"klog/app/cli/lib"
	"klog/service"
	gotime "time"
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

func (args *NowArgs) total(reference gotime.Time, rs ...Record) Duration {
	if args.Now {
		d, _ := service.HypotheticalTotal(reference, rs...)
		return d
	}
	return service.Total(rs...)
}

type FilterArgs struct {
	Tags      []string   `name:"tag" group:"Filter" help:"Only records (or particular entries) that match this tag"`
	Date      []Date     `name:"date" group:"Filter" help:"Only records at this date"`
	Today     bool       `name:"today" group:"Filter" help:"Only records at today’s date"`
	Yesterday bool       `name:"yesterday" group:"Filter" help:"Only records at yesterday’s date"`
	Since     Date       `name:"since" group:"Filter" help:"Only records since this date (inclusive)"`
	Until     Date       `name:"until" group:"Filter" help:"Only records until this date (inclusive)"`
	After     Date       `name:"after" group:"Filter" help:"Only records after this date (exclusive)"`
	Before    Date       `name:"before" group:"Filter" help:"Only records before this date (exclusive)"`
	Period    lib.Period `name:"period" group:"Filter" help:"Only records in this period (YYYY-MM or YYYY)"`
}

func (args *FilterArgs) filter(now gotime.Time, rs []Record) []Record {
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
	return prettifyWarnings(ws)
}

type SortArgs struct {
	Sort string `name:"sort" help:"Sort output by date (ASC or DESC)" enum:"ASC,DESC,"`
}

func (args *SortArgs) sort(rs []Record) []Record {
	if args.Sort == "" {
		return rs
	}
	startWithOldest := false
	if args.Sort == "ASC" {
		startWithOldest = true
	}
	return service.Sort(rs, startWithOldest)
}
