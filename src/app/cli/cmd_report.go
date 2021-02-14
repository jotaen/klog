package cli

import (
	"fmt"
	. "klog"
	"klog/app"
	"klog/service"
	"strings"
	gotime "time"
)

type Report struct {
	DiffArg
	FilterArgs
	WarnArgs
	Fill bool `name:"fill" help:"Show all consecutive days, even if there is no record"`
	NowArgs
	InputFilesArgs
}

func (opt *Report) Run(ctx app.Context) error {
	records, err := ctx.RetrieveRecords(opt.File...)
	if err != nil {
		return err
	}
	if len(records) == 0 {
		return nil
	}
	records = opt.filter(records)
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
	now := gotime.Now()
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
				return prettyMonth(m)[:3]
			}
			return "   "
		}()
		day := func() string {
			return fmt.Sprintf("%s %2v.", prettyDay(date.Weekday())[:3], date.Day())
		}()
		ctx.Print(fmt.Sprintf("%s %s    %s  ", year, month, day))

		rs := recordGroups[date.Hash()]
		if len(rs) == 0 {
			ctx.Print("\n")
			continue
		}
		total := opt.NowArgs.total(now, rs...)
		ctx.Print(pad(7-len(total.ToString())) + styler.Duration(total, false))

		if opt.Diff {
			should := service.ShouldTotalSum(rs...)
			ctx.Print(pad(10-len(should.ToString())) + styler.ShouldTotal(should))
			diff := total.Minus(should)
			ctx.Print(pad(9-len(diff.ToStringWithSign())) + styler.Duration(diff, true))
		}

		ctx.Print("\n")
	}
	ctx.Print(indentation + " " + strings.Repeat("=", 8))
	if opt.Diff {
		ctx.Print(strings.Repeat("=", 19))
	}
	ctx.Print("\n")
	grandTotal := opt.NowArgs.total(now, records...)
	ctx.Print(indentation + pad(9-len(grandTotal.ToStringWithSign())) + styler.Duration(grandTotal, true))
	if opt.Diff {
		grandShould := service.ShouldTotalSum(records...)
		ctx.Print(pad(10-len(grandShould.ToString())) + styler.ShouldTotal(grandShould))
		grandDiff := grandTotal.Minus(grandShould)
		ctx.Print(pad(9-len(grandDiff.ToStringWithSign())) + styler.Duration(grandDiff, true))
	}
	ctx.Print("\n")

	ctx.Print(opt.WarnArgs.ToString(records))
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
