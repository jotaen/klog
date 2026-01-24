package args

import (
	gotime "time"

	"github.com/jotaen/klog/klog"
	"github.com/jotaen/klog/klog/app"
	"github.com/jotaen/klog/klog/service"
	"github.com/jotaen/klog/klog/service/filter"
	"github.com/jotaen/klog/klog/service/period"
)

type FilterArgs struct {
	// General filters
	Date      klog.Date         `name:"date" placeholder:"DATE" group:"Filter" help:"Entries at this date. DATE has to be in format YYYY-MM-DD or YYYY/MM/DD. E.g., '2024-01-31' or '2024/01/31'."`
	Since     klog.Date         `name:"since" placeholder:"DATE" group:"Filter" help:"Entries since this date (inclusive)."`
	Until     klog.Date         `name:"until" placeholder:"DATE" group:"Filter" help:"Entries until this date (inclusive)."`
	Period    period.Period     `name:"period" placeholder:"PERIOD" group:"Filter" help:"Entries within a calendar period. PERIOD has to be in format YYYY, YYYY-MM, YYYY-Www or YYYY-Qq. E.g., '2024', '2024-04', '2022-W21' or '2024-Q1'."`
	Tags      []klog.Tag        `name:"tag" placeholder:"TAG" group:"Filter" help:"Entries that match these tags (either in the record summary or the entry summary). You can omit the leading '#'."`
	EntryType service.EntryType `name:"entry-type" placeholder:"TYPE" group:"Filter" help:"Entries of this type. TYPE can be 'range', 'open-range', 'duration', 'duration-positive' or 'duration-negative'."`

	// Filter shortcuts:
	// The two `XXX` ones are dummy entries just for the help output, they also aren’t available
	// for tab completion. The other ones are not shown in the help output (because that would be
	// too verbose then), but they are still available for tab completion.
	Today       bool `name:"today" group:"Filter" help:"Records at today’s date."`
	Yesterday   bool `name:"yesterday" group:"Filter" help:"Records at yesterday’s date."`
	Tomorrow    bool `name:"tomorrow" group:"Filter" help:"Records at tomorrow’s date."`
	ThisXXX     bool `name:"this-***" group:"Filter" help:"Records of this week/month/quarter/year, e.g. '--this-week' or '--this-quarter'." completion-enabled:"false"`
	LastXXX     bool `name:"last-***" group:"Filter" help:"Records of last week/month/quarter/year, e.g. '--last-month' or '--last-year'." completion-enabled:"false"`
	ThisWeek    bool `hidden:"" name:"this-week" group:"Filter" completion-enabled:"true"`
	LastWeek    bool `hidden:"" name:"last-week" group:"Filter" completion-enabled:"true"`
	ThisMonth   bool `hidden:"" name:"this-month" group:"Filter" completion-enabled:"true"`
	LastMonth   bool `hidden:"" name:"last-month" group:"Filter" completion-enabled:"true"`
	ThisQuarter bool `hidden:"" name:"this-quarter" group:"Filter" completion-enabled:"true"`
	LastQuarter bool `hidden:"" name:"last-quarter" group:"Filter" completion-enabled:"true"`
	ThisYear    bool `hidden:"" name:"this-year" group:"Filter" completion-enabled:"true"`
	LastYear    bool `hidden:"" name:"last-year" group:"Filter" completion-enabled:"true"`

	Filter string `name:"filter" placeholder:"FILTER-EXPRESSION" group:"Filter" help:"Entries that match this filter expression. Run 'klog info --filtering' to learn how to use filter expressions."`
}

func (args *FilterArgs) ApplyFilter(now gotime.Time, rs []klog.Record) ([]klog.Record, app.Error) {
	if args.Filter != "" {
		predicate, err := filter.Parse(args.Filter)
		if err != nil {
			return nil, app.NewErrorWithCode(
				app.GENERAL_ERROR,
				"Malformed filter query",
				err.Error(),
				err,
			)
		}
		rs = filter.Filter(predicate, rs)
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
