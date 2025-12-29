package args

import (
	gotime "time"

	"github.com/jotaen/klog/klog"
	"github.com/jotaen/klog/klog/app"
	"github.com/jotaen/klog/klog/service"
	"github.com/jotaen/klog/klog/service/kfl"
	"github.com/jotaen/klog/klog/service/period"
)

type FilterArgs struct {
	// General filters
	Date      klog.Date         `name:"date" placeholder:"DATE" group:"Filter" help:"Records at this date. DATE has to be in format YYYY-MM-DD or YYYY/MM/DD. E.g., '2024-01-31' or '2024/01/31'."`
	Since     klog.Date         `name:"since" placeholder:"DATE" group:"Filter" help:"Records since this date (inclusive)."`
	Until     klog.Date         `name:"until" placeholder:"DATE" group:"Filter" help:"Records until this date (inclusive)."`
	Period    period.Period     `name:"period" placeholder:"PERIOD" group:"Filter" help:"Records within a calendar period. PERIOD has to be in format YYYY, YYYY-MM, YYYY-Www or YYYY-Qq. E.g., '2024', '2024-04', '2022-W21' or '2024-Q1'."`
	Tags      []klog.Tag        `name:"tag" placeholder:"TAG" group:"Filter" help:"Records (or entries) that match these tags. You can omit the leading '#'."`
	EntryType service.EntryType `name:"entry-type" placeholder:"TYPE" group:"Filter" help:"Entries of this type. TYPE can be 'range', 'open-range', 'duration', 'duration-positive' or 'duration-negative'."`

	// Shortcut filters
	// The `XXX` ones are dummy entries just for the help output
	Today       bool `name:"today" group:"Filter" help:"Records at today’s date."`
	Yesterday   bool `name:"yesterday" group:"Filter" help:"Records at yesterday’s date."`
	Tomorrow    bool `name:"tomorrow" group:"Filter" help:"Records at tomorrow’s date."`
	ThisXXX     bool `name:"this-***" group:"Filter" help:"Records of this week/month/quarter/year, e.g. '--this-week' or '--this-quarter'."`
	LastXXX     bool `name:"last-***" group:"Filter" help:"Records of last week/month/quarter/year, e.g. '--last-month' or '--last-year'."`
	ThisWeek    bool `hidden:"" name:"this-week" group:"Filter"`
	LastWeek    bool `hidden:"" name:"last-week" group:"Filter"`
	ThisMonth   bool `hidden:"" name:"this-month" group:"Filter"`
	LastMonth   bool `hidden:"" name:"last-month" group:"Filter"`
	ThisQuarter bool `hidden:"" name:"this-quarter" group:"Filter"`
	LastQuarter bool `hidden:"" name:"last-quarter" group:"Filter"`
	ThisYear    bool `hidden:"" name:"this-year" group:"Filter"`
	LastYear    bool `hidden:"" name:"last-year" group:"Filter"`

	FilterQuery string `name:"filter" placeholder:"KQL-FILTER-QUERY" group:"Filter" help:"(Experimental)"`
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

func (args *FilterArgs) ApplyFilter(now gotime.Time, rs []klog.Record) ([]klog.Record, app.Error) {
	if args.FilterQuery != "" {
		predicate, err := kfl.Parse(args.FilterQuery)
		if err != nil {
			return nil, app.NewErrorWithCode(
				app.GENERAL_ERROR,
				"Malformed filter query",
				err.Error(),
				err,
			)
		}
		rs = kfl.Filter(predicate, rs)
		return rs, nil
	}
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
	if args.Today {
		qry.AtDate = today
	}
	if args.Yesterday {
		qry.AtDate = today.PlusDays(-1)
	}
	if args.Tomorrow {
		qry.AtDate = today.PlusDays(+1)
	}
	if args.EntryType != "" {
		qry.EntryType = args.EntryType
	}
	shortcutPeriod := func() period.Period {
		if args.ThisWeek {
			return period.NewWeekFromDate(today).Period()
		}
		if args.LastWeek {
			return period.NewWeekFromDate(today).Previous().Period()
		}
		if args.ThisMonth {
			return period.NewMonthFromDate(today).Period()
		}
		if args.LastMonth {
			return period.NewMonthFromDate(today).Previous().Period()
		}
		if args.ThisQuarter {
			return period.NewQuarterFromDate(today).Period()
		}
		if args.LastQuarter {
			return period.NewQuarterFromDate(today).Previous().Period()
		}
		if args.ThisYear {
			return period.NewYearFromDate(today).Period()
		}
		if args.LastYear {
			return period.NewYearFromDate(today).Previous().Period()
		}
		return nil
	}()
	if shortcutPeriod != nil {
		qry.AfterOrEqual = shortcutPeriod.Since()
		qry.BeforeOrEqual = shortcutPeriod.Until()
	}
	return service.Filter(rs, qry), nil
}
