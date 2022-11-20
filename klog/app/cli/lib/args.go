package lib

import (
	"github.com/jotaen/klog/klog"
	"github.com/jotaen/klog/klog/app"
	"github.com/jotaen/klog/klog/service"
	"github.com/jotaen/klog/klog/service/period"
	"strings"
	gotime "time"
)

type InputFilesArgs struct {
	File []app.FileOrBookmarkName `arg:"" optional:"" type:"string" predictor:"file_or_bookmark" name:"file or bookmark" help:".klg source file(s) (if empty the bookmark is used)"`
}

type OutputFileArgs struct {
	File app.FileOrBookmarkName `arg:"" optional:"" type:"string" predictor:"file_or_bookmark" name:"file or bookmark" help:".klg source file (if empty the bookmark is used)"`
}

type AtDateArgs struct {
	Date      klog.Date `name:"date" short:"d" help:"The date of the record"`
	Today     bool      `name:"today" help:"Use today’s date (default)"`
	Yesterday bool      `name:"yesterday" help:"Use yesterday’s date"`
	Tomorrow  bool      `name:"tomorrow" help:"Use tomorrow’s date"`
}

func (args *AtDateArgs) AtDate(now gotime.Time) (klog.Date, bool) {
	if args.Date != nil {
		return args.Date, false
	}
	today := klog.NewDateFromGo(now) // That’s effectively/implicitly `--today`
	if args.Yesterday {
		return today.PlusDays(-1), false
	} else if args.Tomorrow {
		return today.PlusDays(1), false
	}
	return today, true
}

type AtDateAndTimeArgs struct {
	AtDateArgs
	Time  klog.Time        `name:"time" short:"t" help:"Specify the time (defaults to now)"`
	Round service.Rounding `name:"round" short:"r" help:"Round time to nearest multiple of 5m, 10m, 15m, 30m, or 60m / 1h"`
}

func (args *AtDateAndTimeArgs) AtTime(now gotime.Time) (klog.Time, bool, app.Error) {
	if args.Time != nil {
		return args.Time, false, nil
	}
	date, _ := args.AtDate(now)
	today := klog.NewDateFromGo(now)
	time := klog.NewTimeFromGo(now)
	if args.Round != nil {
		time = service.RoundToNearest(time, args.Round)
	}
	if today.IsEqualTo(date) {
		return time, true, nil
	} else if today.PlusDays(-1).IsEqualTo(date) {
		shiftedTime, _ := time.Plus(klog.NewDuration(24, 0))
		return shiftedTime, true, nil
	} else if today.PlusDays(1).IsEqualTo(date) {
		shiftedTime, _ := time.Plus(klog.NewDuration(-24, 0))
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

func (args *NowArgs) ApplyNow(reference gotime.Time, rs ...klog.Record) ([]klog.Record, app.Error) {
	if args.Now {
		rs, err := service.CloseOpenRanges(reference, rs...)
		if err != nil {
			return nil, app.NewErrorWithCode(
				app.LOGICAL_ERROR,
				"Cannot apply --now flag",
				"There are records with uncloseable time ranges",
				err,
			)
		}
		return rs, nil
	}
	return rs, nil
}

type FilterArgs struct {
	// General filters
	Tags   []klog.Tag    `name:"tag" group:"Filter" help:"Records (or entries) that match these tags"`
	Date   klog.Date     `name:"date" group:"Filter" help:"Records at this date"`
	Since  klog.Date     `name:"since" group:"Filter" help:"Records since this date (inclusive)"`
	Until  klog.Date     `name:"until" group:"Filter" help:"Records until this date (inclusive)"`
	After  klog.Date     `name:"after" group:"Filter" help:"Records after this date (exclusive)"`
	Before klog.Date     `name:"before" group:"Filter" help:"Records before this date (exclusive)"`
	Period period.Period `name:"period" group:"Filter" help:"Records in period: YYYY (year), YYYY-MM (month), YYYY-Www (week), or YYYY-Qq (quarter)"`

	// Shortcut filters
	// The `XXX` ones are dummy entries just for the help output
	Today            bool `name:"today" group:"Filter (shortcuts)" help:"Records at today’s date"`
	Yesterday        bool `name:"yesterday" group:"Filter (shortcuts)" help:"Records at yesterday’s date"`
	Tomorrow         bool `name:"tomorrow" group:"Filter (shortcuts)" help:"Records at tomorrow’s date"`
	ThisWeek         bool `name:"this-week" group:"Filter (shortcuts)" help:"Records of the current week"`
	ThisWeekAlias    bool `name:"thisweek" group:"Filter (shortcuts)" hidden:""`
	LastWeek         bool `name:"last-week" group:"Filter (shortcuts)" help:"Records of the last week"`
	LastWeekAlias    bool `name:"lastweek" group:"Filter (shortcuts)" hidden:""`
	ThisMonth        bool `name:"this-month" group:"Filter (shortcuts)" help:"Records of the current month"`
	ThisMonthAlias   bool `name:"thismonth" group:"Filter (shortcuts)" hidden:""`
	LastMonth        bool `name:"last-month" group:"Filter (shortcuts)" help:"Records of the last month"`
	LastMonthAlias   bool `name:"lastmonth" group:"Filter (shortcuts)" hidden:""`
	ThisQuarter      bool `name:"this-quarter" group:"Filter (shortcuts)" help:"Records of the current quarter"`
	ThisQuarterAlias bool `name:"thisquarter" group:"Filter (shortcuts)" hidden:""`
	LastQuarter      bool `name:"last-quarter" group:"Filter (shortcuts)" help:"Records of the last quarter"`
	LastQuarterAlias bool `name:"lastquarter" group:"Filter (shortcuts)" hidden:""`
	ThisYear         bool `name:"this-year" group:"Filter (shortcuts)" help:"Records of the current year"`
	ThisYearAlias    bool `name:"thisyear" group:"Filter (shortcuts)" hidden:""`
	LastYear         bool `name:"last-year" group:"Filter (shortcuts)" help:"Records of the last year"`
	LastYearAlias    bool `name:"lastyear" group:"Filter (shortcuts)" hidden:""`
}

func (args *FilterArgs) ApplyFilter(now gotime.Time, rs []klog.Record) []klog.Record {
	today := klog.NewDateFromGo(now)
	qry := service.FilterQry{
		BeforeOrEqual: args.Until,
		AfterOrEqual:  args.Since,
		Tags:          args.Tags,
		AtDate:        args.Date,
	}
	if args.Period != nil {
		qry.BeforeOrEqual = args.Period.Until()
		qry.AfterOrEqual = args.Period.Since()
	}
	if args.After != nil {
		qry.AfterOrEqual = args.After.PlusDays(1)
	}
	if args.Before != nil {
		qry.BeforeOrEqual = args.Before.PlusDays(-1)
	}
	if args.Today {
		qry.AtDate = today
	}
	if args.Yesterday {
		qry.AtDate = today.PlusDays(-1)
	}
	if args.Tomorrow {
		qry.AtDate = today.PlusDays(+1)
	}
	shortcutPeriod := func() period.Period {
		if args.ThisWeek || args.ThisWeekAlias {
			return period.NewWeekFromDate(today).Period()
		}
		if args.LastWeek || args.LastWeekAlias {
			return period.NewWeekFromDate(today).Previous().Period()
		}
		if args.ThisMonth || args.ThisMonthAlias {
			return period.NewMonthFromDate(today).Period()
		}
		if args.LastMonth || args.LastMonthAlias {
			return period.NewMonthFromDate(today).Previous().Period()
		}
		if args.ThisQuarter || args.ThisQuarterAlias {
			return period.NewQuarterFromDate(today).Period()
		}
		if args.LastQuarter || args.LastQuarterAlias {
			return period.NewQuarterFromDate(today).Previous().Period()
		}
		if args.ThisYear || args.ThisYearAlias {
			return period.NewYearFromDate(today).Period()
		}
		if args.LastYear || args.LastYearAlias {
			return period.NewYearFromDate(today).Previous().Period()
		}
		return nil
	}()
	if shortcutPeriod != nil {
		qry.AfterOrEqual = shortcutPeriod.Since()
		qry.BeforeOrEqual = shortcutPeriod.Until()
	}
	return service.Filter(rs, qry)
}

type WarnArgs struct {
	NoWarn bool `name:"no-warn" help:"Suppress warnings about potential mistakes"`
}

func (args *WarnArgs) PrintWarnings(ctx app.Context, records []klog.Record) {
	if args.NoWarn {
		return
	}
	service.CheckForWarnings(func(w service.Warning) {
		ctx.Print(PrettifyWarning(w))
	}, ctx.Now(), records)
}

type NoStyleArgs struct {
	NoStyle bool `name:"no-style" help:"Do not style or color the values"`
}

func (args *NoStyleArgs) Apply(ctx *app.Context) {
	if args.NoStyle || (*ctx).Preferences().NoColour {
		if s, ok := (*ctx).Serialiser().(CliSerialiser); ok {
			s.Unstyled = true
			(*ctx).SetSerialiser(s)
		}
	}
}

type QuietArgs struct {
	Quiet bool `name:"quiet" help:"Output parseable data without descriptive text"`
}

type SortArgs struct {
	Sort string `name:"sort" help:"Sort output by date (asc or desc)" enum:"asc,desc,ASC,DESC," default:""`
}

func (args *SortArgs) ApplySort(rs []klog.Record) []klog.Record {
	if args.Sort == "" {
		return rs
	}
	startWithOldest := false
	if strings.ToLower(args.Sort) == "asc" {
		startWithOldest = true
	}
	return service.Sort(rs, startWithOldest)
}

type DecimalArgs struct {
	Decimal bool `name:"decimal" help:"Display totals as decimal values (in minutes)"`
}

func (args *DecimalArgs) Apply(ctx *app.Context) {
	if args.Decimal {
		if s, ok := (*ctx).Serialiser().(CliSerialiser); ok {
			s.Decimal = true
			(*ctx).SetSerialiser(s)
		}
	}
}
