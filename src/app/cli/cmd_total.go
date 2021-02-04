package cli

import (
	"errors"
	"fmt"
	"klog"
	"klog/app"
	"klog/service"
	"time"
)

// Deprecated
type Eval struct {
	Total
}

func (args *Eval) Run(ctx *app.Context) error {
	return errors.New("Subcommand `eval` is now named `total`")
}

type Total struct {
	FilterArgs
	MultipleFilesArgs
	Diff bool `name:"diff" help:"Show difference between actual and should total time"`
	Live bool `name:"live" help:"Keep shell open and follow changes live"`
}

func (args *Total) Run(ctx *app.Context) error {
	call := func(f func()) { f() }
	if args.Live {
		call = args.repeat
	}
	call(func() { args.printEvaluation(ctx) })
	return nil
}

func (args *Total) repeat(cb func()) {
	ticker := time.NewTicker(1 * time.Second)
	for time.Now(); true; <-ticker.C {
		fmt.Printf("\033[2J\033[H") // clear screen
		cb()
	}
}

func (args *Total) printEvaluation(ctx *app.Context) {
	rs, err := ctx.RetrieveRecords(args.File...)
	if err != nil {
		fmt.Println(prettifyError(err))
		return
	}
	rs = service.FindFilter(rs, args.toFilter())
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
	fmt.Printf("(In %d record%s)\n", len(rs), func() string {
		if len(rs) == 1 {
			return ""
		}
		return "s"
	}())
}
