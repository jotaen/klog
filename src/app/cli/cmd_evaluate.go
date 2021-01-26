package cli

import (
	"fmt"
	"klog"
	"klog/app"
	"klog/service"
	"time"
)

type Evaluate struct {
	FilterArgs
	MultipleFilesArgs
	Diff bool `name:"diff" help:"Show difference between actual and should total time"`
	Live bool `name:"live" help:"Follow changes in files"`
}

func (args *Evaluate) Run(ctx *app.Context) error {
	call := func(f func()) { f() }
	if args.Live {
		call = args.repeat
	}
	call(func() { args.printEvaluation(ctx) })
	return nil
}

func (args *Evaluate) repeat(cb func()) {
	ticker := time.NewTicker(1 * time.Second)
	for time.Now(); true; <-ticker.C {
		fmt.Printf("\033[2J\033[H") // clear screen
		cb()
	}
}

func (args *Evaluate) printEvaluation(ctx *app.Context) {
	rs, err := ctx.RetrieveRecords(args.File...)
	if err != nil {
		fmt.Println(prettifyError(err))
		return
	}
	rs, _ = service.FindFilter(rs, args.toFilter())
	total, _ := func() (klog.Duration, bool) {
		if args.Live {
			return service.HypotheticalTotal(time.Now(), rs...)
		}
		return service.Total(rs...), false
	}()
	fmt.Printf("Total: %s\n", total.ToString())
	if args.Diff {
		should := service.ShouldTotalSum(rs...)
		diff := klog.NewDuration(0, 0).Minus(should).Plus(total)
		fmt.Printf("Should: %s\n", styler.PrintShouldTotal(should))
		fmt.Printf("Diff: %s\n", styler.PrintDiff(diff))
	}
	fmt.Printf("(In %d records)\n", len(rs))
}
