package cli

import (
	"fmt"
	src "klog"
	"klog/app"
	"klog/service"
	"time"
)

type Follow struct {
	SingleFileArgs
}

func (args *Follow) Run(ctx *app.Context) error {
	ticker := time.NewTicker(1 * time.Second)
	for t := time.Now(); true; t = <-ticker.C {
		r, err := args.getToday(ctx)
		fmt.Printf("\033[2J\033[H") // clear screen
		if err != nil {
			fmt.Println(prettifyError(err))
			continue
		}
		if r == nil {
			fmt.Println("No record for today")
			continue
		}
		currentTotal := service.HypotheticalTotal(r, src.NewTimeFromTime(t))
		fmt.Printf("Date: %s\n", r.Date().ToString())
		fmt.Printf("Total: %s\n", currentTotal.ToString())
	}
	return nil
}

func (args *Follow) getToday(ctx *app.Context) (src.Record, error) {
	rs, err := ctx.RetrieveRecords(args.File)
	if err != nil {
		return nil, err
	}
	date := src.NewDateFromTime(time.Now())
	rs, _ = service.FindFilter(rs, service.Filter{BeforeEq: date, AfterEq: date})
	if len(rs) == 0 {
		return nil, nil
	}
	return rs[0], err
}
