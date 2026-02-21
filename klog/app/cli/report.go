package cli

import (
	"math"
	"strings"

	"github.com/jotaen/klog/klog"
	"github.com/jotaen/klog/klog/app"
	"github.com/jotaen/klog/klog/app/cli/args"
	"github.com/jotaen/klog/klog/app/cli/report"
	"github.com/jotaen/klog/klog/service"
	"github.com/jotaen/klog/klog/service/period"
	tf "github.com/jotaen/klog/lib/terminalformat"
)

type Report struct {
	AggregateBy     string `name:"aggregate" placeholder:"KIND" short:"a" help:"How to aggregate the data. KIND can be 'day' (default), 'week', 'month', 'quarter' or 'year'." enum:"DAY,day,d,WEEK,week,w,MONTH,month,m,QUARTER,quarter,q,YEAR,year,y," default:"day"`
	Fill            bool   `name:"fill" short:"f" help:"Fill any calendar gaps and show a consecutive sequence of dates."`
	Chart           bool   `name:"chart" short:"c" help:"Includes a bar chart rendering, to aid visual comparison."`
	ChartResolution int    `name:"chart-res" help:"Configure the chart resolution. INT must be a positive integer, denoting the minutes per rendered block."`
	args.DiffArgs
	args.FilterArgs
	args.NowArgs
	args.DecimalArgs
	args.WarnArgs
	args.NoStyleArgs
	args.InputFilesArgs
}

func (opt *Report) Help() string {
	return `
It aggregates the totals by period, and prints the respective values chronologically (from oldest to latest).
The default aggregation is by day, but you can choose other periods via the '--aggregate' flag.

The report skips all days (weeks, months, etc.) if no data is available for them.
If you want a consecutive, chronological stream, you can use the '--fill' flag.
`
}

func (opt *Report) Run(ctx app.Context) app.Error {
	opt.DecimalArgs.Apply(&ctx)
	opt.NoStyleArgs.Apply(&ctx)
	cErr := opt.canonicaliseOpts()
	if cErr != nil {
		return cErr
	}
	_, serialiser := ctx.Serialise()
	records, err := ctx.ReadInputs(opt.File...)
	if err != nil {
		return err
	}
	now := ctx.Now()
	records, fErr := opt.ApplyFilter(now, records)
	if fErr != nil {
		return fErr
	}
	if len(records) == 0 {
		return nil
	}
	nErr := opt.ApplyNow(now, records...)
	if nErr != nil {
		return nErr
	}
	records = service.Sort(records, true)
	aggregator := opt.aggregator()
	recordGroups, dates := groupByDate(aggregator.DateHash, records)
	if opt.Fill {
		singlePeriod := opt.FilterArgs.SinglePeriodRequested()
		if singlePeriod != nil {
			dates = allDatesRange(singlePeriod.Since(), singlePeriod.Until())
		} else if len(records) > 0 {
			dates = allDatesRange(records[0].Date(), records[len(records)-1].Date())
		}
	}

	// Table setup
	numberOfValueColumns := func() int {
		n := 1
		if opt.Diff {
			n += 2
		}
		if opt.Chart {
			n += 1
		}
		return n
	}()
	table := tf.NewTable(
		aggregator.NumberOfPrefixColumns()+numberOfValueColumns,
		" ",
	)

	// Header
	aggregator.OnHeaderPrefix(table)
	table.CellR("   Total")
	if opt.Diff {
		table.CellR("   Should").CellR("    Diff")
	}
	if opt.Chart {
		table.Skip(1)
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
		table.CellR(serialiser.Duration(total))

		if opt.Diff {
			should := service.ShouldTotalSum(rs...)
			diff := service.Diff(should, total)
			table.CellR(serialiser.ShouldTotal(should)).CellR(serialiser.SignedDuration(diff))
		}
		if opt.Chart {
			table.CellL(" " + renderBar(opt.ChartResolution, total))
		}
	}

	// Line
	table.Skip(aggregator.NumberOfPrefixColumns()).Fill("=")
	if opt.Diff {
		table.Fill("=").Fill("=")
	}
	if opt.Chart {
		table.Skip(1)
	}

	// Footer
	grandTotal := service.Total(records...)
	table.Skip(aggregator.NumberOfPrefixColumns())
	table.CellR(serialiser.Duration(grandTotal))
	if opt.Diff {
		grandShould := service.ShouldTotalSum(records...)
		grandDiff := service.Diff(grandShould, grandTotal)
		table.CellR(serialiser.ShouldTotal(grandShould)).CellR(serialiser.SignedDuration(grandDiff))
	}
	if opt.Chart {
		table.Skip(1)
	}

	table.Collect(ctx.Print)
	opt.WarnArgs.PrintWarnings(ctx, records, []service.UsageWarning{opt.NowArgs.GetWarning(), opt.DiffArgs.GetWarning(opt.FilterArgs)})
	return nil
}

func (opt *Report) canonicaliseOpts() app.Error {
	if opt.AggregateBy == "" {
		opt.AggregateBy = "d"
	} else {
		opt.AggregateBy = strings.ToLower(opt.AggregateBy[:1])
	}

	if opt.ChartResolution == 0 {
		// If the resolution wasn’t explicitly specified, use a default value
		// that aims for a good balance between granularity and overall row width
		// in the context of the desired aggregation mode.
		switch opt.AggregateBy {
		case "y":
			opt.ChartResolution = 60 * 8 * 7 // Full working week
		case "q":
			opt.ChartResolution = 60 * 8 // Full working day
		case "m":
			opt.ChartResolution = 60 * 4 // Half working day
		case "w":
			opt.ChartResolution = 60
		default: // "d"
			opt.ChartResolution = 15
		}
	} else if opt.ChartResolution > 0 {
		// When chart resolution is specified, automatically assume --chart
		// to be given as well.
		opt.Chart = true
	} else if opt.ChartResolution < 0 {
		return app.NewErrorWithCode(app.LOGICAL_ERROR, "Invalid resolution", "The resolution must be a positive integer", nil)
	}
	return nil
}

func (opt *Report) aggregator() report.Aggregator {
	switch opt.AggregateBy {
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

func renderBar(minutesPerUnit int, d klog.Duration) string {
	block := "▇"
	blocksCount := func() int {
		mins := d.InMinutes()
		if mins <= 0 {
			return 0
		}
		return int(math.Ceil(float64(mins) / float64(minutesPerUnit)))
	}()
	return strings.Repeat(block, blocksCount)
}
