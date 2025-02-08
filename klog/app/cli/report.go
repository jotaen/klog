package cli

import (
	"github.com/jotaen/klog/klog"
	"github.com/jotaen/klog/klog/app"
	"github.com/jotaen/klog/klog/app/cli/report"
	tf "github.com/jotaen/klog/klog/app/cli/terminalformat"
	"github.com/jotaen/klog/klog/app/cli/util"
	"github.com/jotaen/klog/klog/service"
	"github.com/jotaen/klog/klog/service/period"
	"math"
	"strings"
)

type Report struct {
	AggregateBy     string `name:"aggregate" placeholder:"KIND" short:"a" help:"How to aggregate the data. KIND can be 'day' (default), 'week', 'month', 'quarter' or 'year'." enum:"DAY,day,d,WEEK,week,w,MONTH,month,m,QUARTER,quarter,q,YEAR,year,y," default:"day"`
	Fill            bool   `name:"fill" short:"f" help:"Fill any calendar gaps and show a consecutive sequence of dates."`
	Chart           bool   `name:"chart" short:"c" help:"Includes a bar chart rendering, to aid visual comparison."`
	ChartResolution int    `name:"chart-res" help:"Configure the chart resolution. INT must be a positive integer, denoting the minutes per rendered block. The default is 15."`
	util.DiffArgs
	util.FilterArgs
	util.NowArgs
	util.DecimalArgs
	util.WarnArgs
	util.NoStyleArgs
	util.InputFilesArgs
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
	cErr := opt.ApplyChart()
	if cErr != nil {
		return cErr
	}
	_, serialiser := ctx.Serialise()
	records, err := ctx.ReadInputs(opt.File...)
	if err != nil {
		return err
	}
	now := ctx.Now()
	records = opt.ApplyFilter(now, records)
	if len(records) == 0 {
		return nil
	}
	nErr := opt.ApplyNow(now, records...)
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
	opt.WarnArgs.PrintWarnings(ctx, records, opt.GetNowWarnings())
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

func (opt *Report) ApplyChart() app.Error {
	if opt.ChartResolution == 0 {
		// Unless specified otherwise, set the default resolution to 15 minutes
		// per rendered block. This should make for a good balance between granularity
		// and row width in the context of the (default) daily aggregation mode.
		opt.ChartResolution = 15
	} else if opt.ChartResolution > 0 {
		// When chart resolution is specified, automatically assume --chart
		// to be given as well.
		opt.Chart = true
	} else if opt.ChartResolution < 0 {
		return app.NewErrorWithCode(app.LOGICAL_ERROR, "Invalid scale factor", "The scale factor must be positive integer", nil)
	}
	return nil
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
