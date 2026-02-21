package args

import (
	gotime "time"

	"github.com/jotaen/klog/klog"
	"github.com/jotaen/klog/klog/app"
	"github.com/jotaen/klog/klog/service/filter"
	"github.com/jotaen/klog/klog/service/period"
)

type FilterArgs struct {
	// Date-related filters:
	Date   klog.Date     `name:"date" placeholder:"DATE" group:"Filter" help:"Entries at this date. DATE has to be in format YYYY-MM-DD or YYYY/MM/DD. E.g., '2024-01-31' or '2024/01/31'."`
	Since  klog.Date     `name:"since" placeholder:"DATE" group:"Filter" help:"Entries since this date (inclusive)."`
	Until  klog.Date     `name:"until" placeholder:"DATE" group:"Filter" help:"Entries until this date (inclusive)."`
	Period period.Period `name:"period" placeholder:"PERIOD" group:"Filter" help:"Entries within a calendar period. PERIOD has to be in format YYYY, YYYY-MM, YYYY-Www or YYYY-Qq. E.g., '2024', '2024-04', '2022-W21' or '2024-Q1'."`

	// Filter shortcuts:
	// The two `XXX` ones are dummy entries just for the help output, they also aren’t available
	// for tab completion. The other ones are not shown in the help output (because that would be
	// too verbose then), but they are still available for tab completion.
	Today       bool `name:"today" group:"Filter" help:"Records at today’s date."`
	Yesterday   bool `name:"yesterday" group:"Filter" help:"Records at yesterday’s date."`
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

	// General filters:
	Tags   []klog.Tag `name:"tag" placeholder:"TAG" group:"Filter" help:"Entries that match these tags (either in the record summary or the entry summary). You can omit the leading '#'."`
	Filter string     `name:"filter" placeholder:"EXPR" group:"Filter" help:"Entries that match this filter expression. Run 'klog info --filtering' to learn how expressions works."`

	hasPartialRecordsWithShouldTotal bool // Field only for internal use
}

func (args *FilterArgs) ApplyFilter(now gotime.Time, rs []klog.Record) ([]klog.Record, app.Error) {
	var predicates = []filter.Predicate{}

	// Closed date-range filters:
	dateRanges := func() []period.Period {
		today := klog.NewDateFromGo(now)
		var res []period.Period
		if args.Date != nil {
			res = append(res, period.NewPeriod(args.Date, args.Date))
		}
		if args.Today {
			res = append(res, period.NewPeriod(today, today))
		}
		if args.Yesterday {
			res = append(res, period.NewPeriod(today.PlusDays(-1), today.PlusDays(-1)))
		}
		if args.Period != nil {
			res = append(res, args.Period)
		}
		if args.ThisWeek {
			res = append(res, period.NewWeekFromDate(today).Period())
		}
		if args.LastWeek {
			res = append(res, period.NewWeekFromDate(today).Previous().Period())
		}
		if args.ThisMonth {
			res = append(res, period.NewMonthFromDate(today).Period())
		}
		if args.LastMonth {
			res = append(res, period.NewMonthFromDate(today).Previous().Period())
		}
		if args.ThisQuarter {
			res = append(res, period.NewQuarterFromDate(today).Period())
		}
		if args.LastQuarter {
			res = append(res, period.NewQuarterFromDate(today).Previous().Period())
		}
		if args.ThisYear {
			res = append(res, period.NewYearFromDate(today).Period())
		}
		if args.LastYear {
			res = append(res, period.NewYearFromDate(today).Previous().Period())
		}
		return res
	}()
	for _, d := range dateRanges {
		predicates = append(predicates, filter.IsInDateRange{
			From: d.Since(),
			To:   d.Until(),
		})
	}

	// Open date-range filters:
	if args.Since != nil {
		predicates = append(predicates, filter.IsInDateRange{
			From: args.Since,
		})
	}
	if args.Until != nil {
		predicates = append(predicates, filter.IsInDateRange{
			To: args.Until,
		})
	}

	// Tag filters:
	for _, t := range args.Tags {
		predicates = append(predicates, filter.HasTag{
			Tag: t,
		})
	}

	// Filter expression:
	if args.Filter != "" {
		filterPredicate, err := filter.Parse(args.Filter)
		if err != nil {
			return nil, app.NewErrorWithCode(
				app.GENERAL_ERROR,
				"Malformed filter query",
				err.Error(),
				err,
			)
		}
		predicates = append(predicates, filterPredicate)
	}

	// Apply filters, if applicable:
	if len(predicates) > 0 {
		hasPartialRecordsWithShouldTotal := false
		rs, hasPartialRecordsWithShouldTotal = filter.Filter(filter.And{Predicates: predicates}, rs)
		args.hasPartialRecordsWithShouldTotal = hasPartialRecordsWithShouldTotal

	}
	return rs, nil
}
