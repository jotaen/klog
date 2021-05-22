package cli

import (
	"fmt"
	. "klog"
	"klog/app"
	"klog/app/cli/lib"
	"klog/lib/jotaen/terminalformat"
	"klog/service"
)

type Report struct {
	lib.DiffArgs
	lib.FilterArgs
	lib.WarnArgs
	Fill bool `name:"fill" short:"f" help:"Fill the gaps and show consecutive stream of days"`
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
	numberOfValueColumns := func() int {
		if opt.Diff {
			return 3
		}
		return 1
	}()
	numberOfColumns := 4 + numberOfValueColumns
	table := terminalformat.NewTable(numberOfColumns, " ")
	records = service.Sort(records, true)
	table.
		CellL("    ").   // 2020
		CellL("   ").    // Dec
		CellL("      "). // Sun
		CellR("   ").    // 17.
		CellR("   Total")
	if opt.Diff {
		table.CellR("   Should").CellR("    Diff")
	}
	recordGroups, dates := groupByDate(records)
	y := -1
	m := -1
	if opt.Fill {
		dates = allDatesRange(records[0].Date(), records[len(records)-1].Date())
	}
	for _, date := range dates {
		// Year
		if date.Year() != y {
			m = -1 // force month to be recalculated
			table.CellR(fmt.Sprint(date.Year()))
			y = date.Year()
		} else {
			table.Skip(1)
		}

		// Month
		if date.Month() != m {
			m = date.Month()
			table.CellR(lib.PrettyMonth(m)[:3])
		} else {
			table.Skip(1)
		}

		// Day
		table.CellR(lib.PrettyDay(date.Weekday())[:3]).CellR(fmt.Sprintf("%2v.", date.Day()))

		rs := recordGroups[date.Hash()]
		if len(rs) == 0 {
			table.Skip(numberOfValueColumns)
			continue
		}
		total := opt.NowArgs.Total(now, rs...)
		table.CellR(ctx.Serialiser().Duration(total))
		if opt.Diff {
			should := service.ShouldTotalSum(rs...)
			diff := service.Diff(should, total)
			table.CellR(ctx.Serialiser().ShouldTotal(should)).CellR(ctx.Serialiser().SignedDuration(diff))
		}
	}

	// Line
	table.Skip(4).Fill("=")
	if opt.Diff {
		table.Fill("=").Fill("=")
	}
	ctx.Print("\n")
	grandTotal := opt.NowArgs.Total(now, records...)

	// Totals
	table.Skip(4)
	table.CellR(ctx.Serialiser().Duration(grandTotal))
	if opt.Diff {
		grandShould := service.ShouldTotalSum(records...)
		grandDiff := service.Diff(grandShould, grandTotal)
		table.CellR(ctx.Serialiser().ShouldTotal(grandShould)).CellR(ctx.Serialiser().SignedDuration(grandDiff))
	}

	ctx.Print(table.ToString())
	ctx.Print(opt.WarnArgs.ToString(now, records))
	return nil
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

func groupByDate(rs []Record) (map[DateHash][]Record, []Date) {
	days := make(map[DateHash][]Record, len(rs))
	var order []Date
	for _, r := range rs {
		h := r.Date().Hash()
		if _, ok := days[h]; !ok {
			days[h] = []Record{}
			order = append(order, r.Date())
		}
		days[h] = append(days[h], r)
	}
	return days, order
}
