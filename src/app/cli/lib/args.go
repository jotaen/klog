package lib

import (
	. "github.com/jotaen/klog/src"
	"github.com/jotaen/klog/src/app"
	"github.com/jotaen/klog/src/parser"
	"github.com/jotaen/klog/src/service"
	"github.com/jotaen/klog/src/service/period"
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
	Time  Time             `name:"time" short:"t" help:"Specify the time (defaults to now)"`
	Round service.Rounding `name:"round" short:"r" help:"Round time to nearest multiple of 5m, 10m, 15m, 30m, or 60m / 1h"`
}

func (args *AtDateAndTimeArgs) AtTime(now gotime.Time) (Time, bool, app.Error) {
	if args.Time != nil {
		return args.Time, false, nil
	}
	date, _ := args.AtDate(now)
	today := NewDateFromGo(now)
	time := NewTimeFromGo(now)
	if args.Round != nil {
		time = service.RoundToNearest(time, args.Round)
	}
	if today.IsEqualTo(date) {
		return time, true, nil
	} else if today.PlusDays(-1).IsEqualTo(date) {
		shiftedTime, _ := time.Plus(NewDuration(24, 0))
		return shiftedTime, true, nil
	} else if today.PlusDays(1).IsEqualTo(date) {
		shiftedTime, _ := time.Plus(NewDuration(-24, 0))
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
	// General filters
	Tags   []Tag         `name:"tag" group:"Filter" help:"Records (or entries) that match these tags"`
	Date   Date          `name:"date" group:"Filter" help:"Records at this date"`
	Since  Date          `name:"since" group:"Filter" help:"Records since this date (inclusive)"`
	Until  Date          `name:"until" group:"Filter" help:"Records until this date (inclusive)"`
	After  Date          `name:"after" group:"Filter" help:"Records after this date (exclusive)"`
	Before Date          `name:"before" group:"Filter" help:"Records before this date (exclusive)"`
	Period period.Period `name:"period" group:"Filter" help:"Records in period: YYYY (year), YYYY-MM (month), YYYY-Www (week), or YYYY-Qq (quarter)"`

	// Shortcut filters
	// The `XXX` ones are dummy entries just for the help output
	Today            bool `name:"today" group:"Filter (shortcuts)" help:"Records at today’s date"`
	Yesterday        bool `name:"yesterday" group:"Filter (shortcuts)" help:"Records at yesterday’s date"`
	Tomorrow         bool `name:"tomorrow" group:"Filter (shortcuts)" help:"Records at tomorrow’s date"`
	ThisXXX          bool `name:"this-***" group:"Filter (shortcuts)" help:"Records of the current week/quarter/month/year (e.g. --this-year)"`
	LastXXX          bool `name:"last-***" group:"Filter (shortcuts)" help:"Records of the previous week/quarter/month/year (e.g. --last-month)"`
	ThisWeek         bool `name:"this-week" hidden:""`
	ThisWeekAlias    bool `name:"thisweek" hidden:""`
	LastWeek         bool `name:"last-week" hidden:""`
	LastWeekAlias    bool `name:"lastweek" hidden:""`
	ThisMonth        bool `name:"this-month" hidden:""`
	ThisMonthAlias   bool `name:"thismonth" hidden:""`
	LastMonth        bool `name:"last-month" hidden:""`
	LastMonthAlias   bool `name:"lastmonth" hidden:""`
	ThisQuarter      bool `name:"this-quarter" hidden:""`
	ThisQuarterAlias bool `name:"thisquarter" hidden:""`
	LastQuarter      bool `name:"last-quarter" hidden:""`
	LastQuarterAlias bool `name:"lastquarter" hidden:""`
	ThisYear         bool `name:"this-year" hidden:""`
	ThisYearAlias    bool `name:"thisyear" hidden:""`
	LastYear         bool `name:"last-year" hidden:""`
	LastYearAlias    bool `name:"lastyear" hidden:""`

	Forth bool `name:"forth" group:"Forth filters"` // TODO write help text
	And   bool `name:"and" group:"Forth filters"`
	Or    bool `name:"or" group:"Forth filters"`
	Not   bool `name:"not" group:"Forth filters"`
}

func flagWithValue(args []string) (string, string) {
	if len(args) == 0 {
		return "", ""
	}
	if len(args) == 1 {
		return args[0], ""
	}
	return args[0], args[1]
}

func applyForthFilter(now gotime.Time, rs []Record) []Record {
	var stack []service.Matcher
	args := os.Args[3:]

	for len(args) > 0 {
		flag, value := flagWithValue(args)
		switch flag {
		case "--date":
			date, _ := NewDateFromString(value)
			matcher := service.DateMatcher(date)
			stack = append(stack, matcher)
			args = args[2:]
		case "--tag":
			tag, _ := NewTagFromString(value)
			matcher := service.TagMatcher(tag)
			stack = append(stack, matcher)
			args = args[2:]
		case "--and":
			if len(stack) >= 2 {
				stack = append(stack[0:len(stack)-2], service.AndMatcher(stack[len(stack)-2], stack[len(stack)-1]))
			}
			args = args[1:]
		case "--or":
			if len(stack) >= 2 {
				stack = append(stack[0:len(stack)-2], service.OrMatcher(stack[len(stack)-2], stack[len(stack)-1]))
			}
			args = args[1:]
		case "--not":
			if len(stack) >= 1 {
				stack = append(stack[0:len(stack)-1], service.NotMatcher(stack[len(stack)-1]))
			}
			args = args[1:]
		default:
			args = args[1:]
		}
	}

	matcher := service.IdentityMatcher
	if len(stack) == 1 {
		matcher = stack[0]
	}
	return service.ForthFilter(matcher, rs)
}

func (args *FilterArgs) ApplyFilter(now gotime.Time, rs []Record) []Record {
	today := NewDateFromGo(now)
	if args.Forth {
		return applyForthFilter(now, rs)
	}
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

func (args *WarnArgs) PrintWarnings(ctx app.Context, records []Record) {
	if args.NoWarn {
		return
	}
	ws := service.CheckForWarnings(ctx.Now(), records)
	ctx.Print(PrettifyWarnings(ws))
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
