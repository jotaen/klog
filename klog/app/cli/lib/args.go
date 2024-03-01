package lib

import (
	"github.com/jotaen/klog/klog"
	"github.com/jotaen/klog/klog/app"
	tf "github.com/jotaen/klog/klog/app/cli/lib/terminalformat"
	"github.com/jotaen/klog/klog/parser/reconciling"
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

func (args *AtDateArgs) AtDate(now gotime.Time) klog.Date {
	if args.Date != nil {
		return args.Date
	}
	today := klog.NewDateFromGo(now) // That’s effectively/implicitly `--today`
	if args.Yesterday {
		return today.PlusDays(-1)
	} else if args.Tomorrow {
		return today.PlusDays(1)
	}
	return today
}

func (args *AtDateArgs) DateFormat(config app.Config) reconciling.ReformatDirective[klog.DateFormat] {
	if args.Date != nil {
		return reconciling.NoReformat[klog.DateFormat]()
	}
	fd := reconciling.ReformatAutoStyle[klog.DateFormat]()
	config.DateUseDashes.Unwrap(func(x bool) {
		fd = reconciling.ReformatExplicitly(klog.DateFormat{UseDashes: x})
	})
	return fd
}

type AtDateAndTimeArgs struct {
	AtDateArgs
	Time  klog.Time        `name:"time" short:"t" help:"Specify the time (defaults to now)"`
	Round service.Rounding `name:"round" short:"r" help:"Round time to nearest multiple of 5m, 10m, 12m, 15m, 20m, 30m, or 60m / 1h"`
}

func (args *AtDateAndTimeArgs) AtTime(now gotime.Time, config app.Config) (klog.Time, app.Error) {
	if args.Time != nil {
		return args.Time, nil
	}
	date := args.AtDate(now)
	today := klog.NewDateFromGo(now)
	time := klog.NewTimeFromGo(now)
	if args.Round != nil {
		time = service.RoundToNearest(time, args.Round)
	} else {
		config.DefaultRounding.Unwrap(func(r service.Rounding) {
			time = service.RoundToNearest(time, r)
		})
	}
	if today.IsEqualTo(date) {
		return time, nil
	} else if today.PlusDays(-1).IsEqualTo(date) {
		shiftedTime, _ := time.Plus(klog.NewDuration(24, 0))
		return shiftedTime, nil
	} else if today.PlusDays(1).IsEqualTo(date) {
		shiftedTime, _ := time.Plus(klog.NewDuration(-24, 0))
		return shiftedTime, nil
	}
	return nil, app.NewErrorWithCode(
		app.LOGICAL_ERROR,
		"Missing time parameter",
		"Please specify a time value for dates in the past",
		nil,
	)
}

func (args *AtDateAndTimeArgs) TimeFormat(config app.Config) reconciling.ReformatDirective[klog.TimeFormat] {
	if args.Time != nil {
		return reconciling.NoReformat[klog.TimeFormat]()
	}
	fd := reconciling.ReformatAutoStyle[klog.TimeFormat]()
	config.TimeUse24HourClock.Unwrap(func(x bool) {
		fd = reconciling.ReformatExplicitly(klog.TimeFormat{Use24HourClock: x})
	})
	return fd
}

func (args *AtDateAndTimeArgs) WasAutomatic() bool {
	return args.Date == nil && args.Time == nil
}

type DiffArgs struct {
	Diff bool `name:"diff" short:"d" help:"Show difference between actual and should-total time"`
}

type NowArgs struct {
	Now          bool `name:"now" short:"n" help:"Assume open ranges to be closed at this moment"`
	hadOpenRange bool // Field only for internal use
}

func (args *NowArgs) ApplyNow(reference gotime.Time, rs ...klog.Record) app.Error {
	if args.Now {
		hasClosedAnyRange, err := service.CloseOpenRanges(reference, rs...)
		if err != nil {
			return app.NewErrorWithCode(
				app.LOGICAL_ERROR,
				"Cannot apply --now flag",
				"There are records with uncloseable time ranges",
				err,
			)
		}
		args.hadOpenRange = hasClosedAnyRange
		return nil
	}
	return nil
}

func (args *NowArgs) HadOpenRange() bool {
	return args.hadOpenRange
}

func (args *NowArgs) GetNowWarnings() []string {
	if args.Now && !args.hadOpenRange {
		return []string{"You specified --now, but there was no open-ended time range."}
	}
	return nil
}

type FilterArgs struct {
	// General filters
	Tags   []klog.Tag    `name:"tag" group:"Filter" help:"Records (or entries) that match these tags"`
	Date   klog.Date     `name:"date" group:"Filter" help:"Records at this date"`
	Since  klog.Date     `name:"since" group:"Filter" help:"Records since this date (inclusive)"`
	Until  klog.Date     `name:"until" group:"Filter" help:"Records until this date (inclusive)"`
	After  klog.Date     `name:"after" group:"Filter" help:"Records after this date (exclusive)"`
	Before klog.Date     `name:"before" group:"Filter" help:"Records before this date (exclusive)"`
	Period period.Period `name:"period" group:"Filter" help:"Records in period: YYYY, YYYY-MM, YYYY-Www, or YYYY-Qq"`

	// Shortcut filters
	// The `XXX` ones are dummy entries just for the help output
	Today            bool `name:"today" group:"Filter" help:"Records at today’s date"`
	Yesterday        bool `name:"yesterday" group:"Filter" help:"Records at yesterday’s date"`
	Tomorrow         bool `name:"tomorrow" group:"Filter" help:"Records at tomorrow’s date"`
	ThisXXX          bool `name:"this-***" group:"Filter" help:"Records of this week/month/quarter/year, e.g. --this-week"`
	LastXXX          bool `name:"last-***" group:"Filter" help:"Records of last week/month/quarter/year, e.g. --last-month"`
	ThisWeek         bool `hidden:"" name:"this-week" group:"Filter"`
	ThisWeekAlias    bool `hidden:"" name:"thisweek" group:"Filter"`
	LastWeek         bool `hidden:"" name:"last-week" group:"Filter"`
	LastWeekAlias    bool `hidden:"" name:"lastweek" group:"Filter"`
	ThisMonth        bool `hidden:"" name:"this-month" group:"Filter"`
	ThisMonthAlias   bool `hidden:"" name:"thismonth" group:"Filter"`
	LastMonth        bool `hidden:"" name:"last-month" group:"Filter"`
	LastMonthAlias   bool `hidden:"" name:"lastmonth" group:"Filter"`
	ThisQuarter      bool `hidden:"" name:"this-quarter" group:"Filter"`
	ThisQuarterAlias bool `hidden:"" name:"thisquarter" group:"Filter"`
	LastQuarter      bool `hidden:"" name:"last-quarter" group:"Filter"`
	LastQuarterAlias bool `hidden:"" name:"lastquarter" group:"Filter"`
	ThisYear         bool `hidden:"" name:"this-year" group:"Filter"`
	ThisYearAlias    bool `hidden:"" name:"thisyear" group:"Filter"`
	LastYear         bool `hidden:"" name:"last-year" group:"Filter"`
	LastYearAlias    bool `hidden:"" name:"lastyear" group:"Filter"`
}

// FilterArgsCompletionOverrides enables/disables tab completion for
// certain flags.
var FilterArgsCompletionOverrides = map[string]bool{
	"this-***":     false, // disable, although not flagged as hidden
	"last-***":     false,
	"this-week":    true, // enable, although flagged as hidden
	"last-week":    true,
	"this-month":   true,
	"last-month":   true,
	"this-quarter": true,
	"last-quarter": true,
	"this-year":    true,
	"last-year":    true,
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

func (args *WarnArgs) PrintWarnings(ctx app.Context, records []klog.Record, additionalWarnings []string) {
	styler, _ := ctx.Serialise()
	if args.NoWarn {
		return
	}
	for _, msg := range additionalWarnings {
		ctx.Print(PrettifyGeneralWarning(msg, styler))
	}
	service.CheckForWarnings(func(w service.Warning) {
		ctx.Print(PrettifyWarning(w, styler))
	}, ctx.Now(), records)
}

type NoStyleArgs struct {
	NoStyle bool `name:"no-style" help:"Do not style or color the values"`
}

func (args *NoStyleArgs) Apply(ctx *app.Context) {
	if args.NoStyle {
		(*ctx).ConfigureSerialisation(func(styler tf.Styler, decimalDuration bool) (tf.Styler, bool) {
			return tf.NewStyler(tf.NO_COLOUR), decimalDuration
		})
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
		(*ctx).ConfigureSerialisation(func(styler tf.Styler, decimalDuration bool) (tf.Styler, bool) {
			return styler, true
		})
	}
}
