package cli

import (
	"fmt"
	"klog/app"
	"klog/service"
	"strings"
)

type Report struct {
	MultipleFilesArgs
}

func (args *Report) Run(ctx app.Context) error {
	rs, err := ctx.RetrieveRecords(args.File...)
	if err != nil {
		return prettifyError(err)
	}
	if len(rs) == 0 {
		return nil
	}
	rs = service.Sort(rs, false)
	fmt.Printf("%s Total\n", strings.Repeat(" ", 22))
	fmt.Printf("%s\n", strings.Repeat("-", 28))
	y := -1
	m := -1
	for _, r := range rs {
		year := func() string {
			if r.Date().Year() != y {
				y = r.Date().Year()
				m = -1 // force month to be recalculated
				return fmt.Sprintf("%d", y)
			}
			return "    "
		}()
		month := func() string {
			if r.Date().Month() != m {
				m = r.Date().Month()
				return prettyMonth(m)
			}
			return "   "
		}()
		day := func() string {
			return fmt.Sprintf("%s %2v.", prettyDay(r.Date().Weekday()), r.Date().Day())
		}()
		fmt.Printf(
			"%s %s    %s    %6v\n",
			year, month, day, service.Total(r).ToString(),
		)
	}
	fmt.Printf("%s\n", strings.Repeat("-", 28))
	fmt.Printf("%28v\n", service.Total(rs...).ToString())
	return nil
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
