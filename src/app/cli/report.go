package cli

import (
	. "klog"
	"klog/app"
	"klog/app/cli/aggregators"
	"klog/app/cli/lib"
	"klog/lib/jotaen/terminalformat"
	"klog/service"
)

type Report struct {
	AggregateBy string `name:"by" default:"day" help:"Aggregate by different categories" enum:"DAY,day,WEEK,week"`
	Fill        bool   `name:"fill" short:"f" help:"Fill the gaps and show consecutive stream of days"`
	lib.DiffArgs
	lib.FilterArgs
	lib.WarnArgs
	lib.NowArgs
	lib.NoStyleArgs
	lib.InputFilesArgs
}

func (opt *Report) Run(ctx app.Context) error {
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
	rowsSeen := make(map[report.Hash]bool)
	for _, date := range dates {
		hash := aggregator.DateHash(date)
		if rowsSeen[hash] {
			continue
		}
		rowsSeen[hash] = true
		aggregator.OnRowPrefix(table, date)
		rs := recordGroups[hash]
		if len(rs) == 0 {
			table.Skip(numberOfValueColumns)
			continue
		}
		// Total
		total := opt.NowArgs.Total(now, rs...)
		table.CellR(ctx.Serialiser().Duration(total))

		// Should/Diff
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
	grandTotal := opt.NowArgs.Total(now, records...)

	// Footer
	table.Skip(aggregator.NumberOfPrefixColumns())
	table.CellR(ctx.Serialiser().Duration(grandTotal))
	if opt.Diff {
		grandShould := service.ShouldTotalSum(records...)
		grandDiff := service.Diff(grandShould, grandTotal)
		table.CellR(ctx.Serialiser().ShouldTotal(grandShould)).CellR(ctx.Serialiser().SignedDuration(grandDiff))
	}

	table.Collect(ctx.Print)
	ctx.Print(opt.WarnArgs.ToString(now, records))
	return nil
}

func (opt *Report) findAggregator() report.Aggregator {
	switch opt.AggregateBy {
	case "week":
		return report.NewWeekAggregator()
	}
	return report.NewDayAggregator()
}

func allDatesRange(from Date, to Date) []Date {
	result := []Date{from}
	for true {
		last := result[len(result)-1]
		if last.IsAfterOrEqual(to) {
			break
		}
		result = append(result, last.PlusDays(1))
	}
	return result
}

func groupByDate(hashProvider func(Date) report.Hash, rs []Record) (map[report.Hash][]Record, []Date) {
	days := make(map[report.Hash][]Record, len(rs))
	var order []Date
	for _, r := range rs {
		h := hashProvider(r.Date())
		if _, ok := days[h]; !ok {
			days[h] = []Record{}
			order = append(order, r.Date())
		}
		days[h] = append(days[h], r)
	}
	return days, order
}
