package cli

import (
	"fmt"
	. "klog"
	"klog/app"
	"klog/app/cli/lib"
	"klog/service"
	"strings"
)

type Report struct {
	lib.DiffArgs
	lib.FilterArgs
	lib.WarnArgs
	Fill bool `name:"fill" short:"f" help:"Show all consecutive days, even if there is no record"`
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
	indentation := strings.Repeat(" ", len("2020 Dec   Wed 30. "))
	records = service.Sort(records, true)
	ctx.Print(indentation + "    Total")
	if opt.Diff {
		ctx.Print("    Should     Diff")
	}
	ctx.Print("\n")
	y := -1
	m := -1
	recordGroups, dates := groupByDate(records)
	if opt.Fill {
		dates = allDatesRange(records[0].Date(), records[len(records)-1].Date())
	}
	for _, date := range dates {
		year := func() string {
			if date.Year() != y {
				y = date.Year()
				m = -1 // force month to be recalculated
				return fmt.Sprintf("%d", y)
			}
			return "    "
		}()
		month := func() string {
			if date.Month() != m {
				m = date.Month()
				return lib.PrettyMonth(m)[:3]
			}
			return "   "
		}()
		day := func() string {
			return fmt.Sprintf("%s %2v.", lib.PrettyDay(date.Weekday())[:3], date.Day())
		}()
		ctx.Print(fmt.Sprintf("%s %s    %s  ", year, month, day))

		rs := recordGroups[date.Hash()]
		if len(rs) == 0 {
			ctx.Print("\n")
			continue
		}
		total := opt.NowArgs.Total(now, rs...)
		ctx.Print(lib.Pad(7-len(total.ToString())) + ctx.Serialiser().Duration(total))

		if opt.Diff {
			should := service.ShouldTotalSum(rs...)
			ctx.Print(lib.Pad(10-len(should.ToString())) + ctx.Serialiser().ShouldTotal(should))
			diff := service.Diff(should, total)
			ctx.Print(lib.Pad(9-len(diff.ToStringWithSign())) + ctx.Serialiser().SignedDuration(diff))
		}

		ctx.Print("\n")
	}
	ctx.Print(indentation + " " + strings.Repeat("=", 8))
	if opt.Diff {
		ctx.Print(strings.Repeat("=", 19))
	}
	ctx.Print("\n")
	grandTotal := opt.NowArgs.Total(now, records...)
	ctx.Print(indentation + lib.Pad(9-len(grandTotal.ToStringWithSign())) + ctx.Serialiser().SignedDuration(grandTotal))
	if opt.Diff {
		grandShould := service.ShouldTotalSum(records...)
		ctx.Print(lib.Pad(10-len(grandShould.ToString())) + ctx.Serialiser().ShouldTotal(grandShould))
		grandDiff := service.Diff(grandShould, grandTotal)
		ctx.Print(lib.Pad(9-len(grandDiff.ToStringWithSign())) + ctx.Serialiser().SignedDuration(grandDiff))
	}
	ctx.Print("\n")

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
