package cli

import (
	"github.com/jotaen/klog/klog"
	"github.com/jotaen/klog/klog/app"
	"github.com/jotaen/klog/klog/app/cli/lib"
	"github.com/jotaen/klog/klog/app/cli/lib/terminalformat"
	"github.com/jotaen/klog/klog/app/cli/report"
	"github.com/jotaen/klog/klog/service"
	"github.com/jotaen/klog/klog/service/period"
	"strings"
)

type Report struct {
	AggregateBy string `name:"aggregate" short:"a" help:"Aggregate data by: day, week, month, quarter, year" enum:"DAY,day,d,WEEK,week,w,MONTH,month,m,QUARTER,quarter,q,YEAR,year,y," default:"day"`
	Fill        bool   `name:"fill" short:"f" help:"Fill the gaps and show a consecutive stream"`
	lib.DiffArgs
	lib.FilterArgs
	lib.NowArgs
	lib.DecimalArgs
	lib.WarnArgs
	lib.NoStyleArgs
	lib.InputFilesArgs
}

func (opt *Report) Run(ctx app.Context) error {
	opt.DecimalArgs.Apply(&ctx)
	opt.NoStyleArgs.Apply(&ctx)
	records, err := ctx.ReadInputs(opt.File...)
	if err != nil {
		return err
	}
	if len(records) == 0 {
		return nil
	}
	now := ctx.Now()
	records = opt.ApplyFilter(now, records)
	records, nErr := opt.ApplyNow(now, records...)
	if nErr != nil {
		return nErr
	}
	records = service.Sort(records, true)
	aggregator := opt.findAggregator()
	recordGroups, dates := groupByDate(aggregator.DateHash, records)
	if opt.Fill {
		dates = allDatesRange(records[0].Date(), records[len(records)-1].Date())
	}

	// Table setup
	numberOfValueColumns := func() int {
		if opt.Diff {
			return 3
		}
		return 1
	}()
	table := terminalformat.NewTable(
		aggregator.NumberOfPrefixColumns()+numberOfValueColumns,
		" ",
	)

	// Header
	aggregator.OnHeaderPrefix(table)
	table.CellR("   Total")
	if opt.Diff {
		table.CellR("   Should").CellR("    Diff")
	}

	// Rows
	hashesAlreadyProcessed := make(map[period.Hash]bool)
	for _, date := range dates {
		hash := aggregator.DateHash(date)
		if hashesAlreadyProcessed[hash] {
			continue
		}
		hashesAlreadyProcessed[hash] = true
		aggregator.OnRowPrefix(table, date)
		rs := recordGroups[hash]
		if len(rs) == 0 {
			table.Skip(numberOfValueColumns)
			continue
		}

		total := service.Total(rs...)
		table.CellR(ctx.Serialiser().Duration(total))

		if opt.Diff {
			should := service.ShouldTotalSum(rs...)
			diff := service.Diff(should, total)
			table.CellR(ctx.Serialiser().ShouldTotal(should)).CellR(ctx.Serialiser().SignedDuration(diff))
		}
	}

	// Line
	table.Skip(aggregator.NumberOfPrefixColumns()).Fill("=")
	if opt.Diff {
		table.Fill("=").Fill("=")
	}
	ctx.Print("\n")
	grandTotal := service.Total(records...)

	// Footer
	table.Skip(aggregator.NumberOfPrefixColumns())
	table.CellR(ctx.Serialiser().Duration(grandTotal))
	if opt.Diff {
		grandShould := service.ShouldTotalSum(records...)
		grandDiff := service.Diff(grandShould, grandTotal)
		table.CellR(ctx.Serialiser().ShouldTotal(grandShould)).CellR(ctx.Serialiser().SignedDuration(grandDiff))
	}

	table.Collect(ctx.Print)
	opt.WarnArgs.PrintWarnings(ctx, records)
	return nil
}

func (opt *Report) findAggregator() report.Aggregator {
	category := (func() string {
		if opt.AggregateBy == "" {
			return "d"
		} else {
			return strings.ToLower(opt.AggregateBy[:1])
		}
	})()
	switch category {
	case "y":
		return report.NewYearAggregator()
	case "q":
		return report.NewQuarterAggregator()
	case "m":
		return report.NewMonthAggregator()
	case "w":
		return report.NewWeekAggregator()
	default: // "d"
		return report.NewDayAggregator()
	}
}

func allDatesRange(from klog.Date, to klog.Date) []klog.Date {
	result := []klog.Date{from}
	for {
		last := result[len(result)-1]
		if last.IsAfterOrEqual(to) {
			break
		}
		result = append(result, last.PlusDays(1))
	}
	return result
}

func groupByDate(hashProvider func(klog.Date) period.Hash, rs []klog.Record) (map[period.Hash][]klog.Record, []klog.Date) {
	days := make(map[period.Hash][]klog.Record, len(rs))
	var order []klog.Date
	for _, r := range rs {
		h := hashProvider(r.Date())
		if _, ok := days[h]; !ok {
			days[h] = []klog.Record{}
			order = append(order, r.Date())
		}
		days[h] = append(days[h], r)
	}
	return days, order
}
