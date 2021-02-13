package cli

import (
	"fmt"
	. "klog"
	"klog/app"
	"klog/service"
	"strings"
)

type Report struct {
	DiffArg
	FilterArgs
	Fill bool `name:"fill" help:"Show all consecutive days, even if there is no record"`
	InputFilesArgs
}

func (args *Report) Run(ctx app.Context) error {
	records, err := ctx.RetrieveRecords(args.File...)
	if err != nil {
		return err
	}
	if len(records) == 0 {
		return nil
	}
	opts := args.FilterArgs.toFilter()
	records = service.Query(records, opts)
	indentation := strings.Repeat(" ", len("2020 Dec   We 30. "))
	records = service.Query(records, service.Opts{Sort: "ASC"})
	ctx.Print(indentation + "    Total")
	if args.Diff {
		ctx.Print("    Should     Diff")
	}
	ctx.Print("\n")
	y := -1
	m := -1
	recordGroups, dates := groupByDate(records)
	if args.Fill {
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
				return prettyMonth(m)
			}
			return "   "
		}()
		day := func() string {
			return fmt.Sprintf("%s %2v.", prettyDay(date.Weekday()), date.Day())
		}()
		ctx.Print(fmt.Sprintf("%s %s    %s  ", year, month, day))

		rs := recordGroups[date.Hash()]
		if len(rs) == 0 {
			ctx.Print("\n")
			continue
		}
		total := service.Total(rs...)
		ctx.Print(pad(7-len(total.ToString())) + styler.Duration(total, false))

		if args.Diff {
			should := service.ShouldTotalSum(rs...)
			ctx.Print(pad(10-len(should.ToString())) + styler.ShouldTotal(should))
			diff := total.Minus(should)
			ctx.Print(pad(9-len(diff.ToStringWithSign())) + styler.Duration(diff, true))
		}

		ctx.Print("\n")
	}
	ctx.Print(indentation + " " + strings.Repeat("=", 8))
	if args.Diff {
		ctx.Print(strings.Repeat("=", 19))
	}
	ctx.Print("\n")
	grandTotal := service.Total(records...)
	ctx.Print(indentation + pad(9-len(grandTotal.ToStringWithSign())) + styler.Duration(grandTotal, true))
	if args.Diff {
		grandShould := service.ShouldTotalSum(records...)
		ctx.Print(pad(10-len(grandShould.ToString())) + styler.ShouldTotal(grandShould))
		grandDiff := grandTotal.Minus(grandShould)
		ctx.Print(pad(9-len(grandDiff.ToStringWithSign())) + styler.Duration(grandDiff, true))
	}
	ctx.Print("\n")
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

func pad(length int) string {
	if length < 0 {
		return ""
	}
	return strings.Repeat(" ", length)
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

func prettyMonth(m int) string {
	switch m {
	case 1:
		return "Jan"
	case 2:
		return "Feb"
	case 3:
		return "Mar"
	case 4:
		return "Apr"
	case 5:
		return "May"
	case 6:
		return "Jun"
	case 7:
		return "Jul"
	case 8:
		return "Aug"
	case 9:
		return "Sep"
	case 10:
		return "Oct"
	case 11:
		return "Nov"
	case 12:
		return "Dec"
	}
	panic("Illegal month") // this can/should never happen
}

func prettyDay(d int) string {
	switch d {
	case 1:
		return "Mo"
	case 2:
		return "Tu"
	case 3:
		return "We"
	case 4:
		return "Th"
	case 5:
		return "Fr"
	case 6:
		return "Sa"
	case 7:
		return "Su"
	}
	panic("Illegal weekday") // this can/should never happen
}
