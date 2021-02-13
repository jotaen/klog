package cli

import (
	"errors"
	"fmt"
	. "klog"
	"klog/app"
	"klog/service"
	"time"
)

// Deprecated
type Eval struct {
	Total
}

func (args *Eval) Run(_ app.Context) error {
	return errors.New("Subcommand `eval` is now named `total`")
}

type Total struct {
	FilterArgs
	DiffArg
	WarnArgs
	Live bool `name:"live" help:"Keep shell open and follow changes live"`
	InputFilesArgs
}

func (args *Total) Run(ctx app.Context) error {
	call := func(f func(ctx app.Context) error) error { return f(ctx) }
	if args.Live {
		call = func(f func(ctx app.Context) error) error { return args.repeat(ctx, f) }
	}
	return call(args.printEvaluation)
}

func (args *Total) repeat(ctx app.Context, cb func(ctx app.Context) error) error {
	ticker := time.NewTicker(1 * time.Second)
	for time.Now(); true; <-ticker.C {
		ctx.Print(fmt.Sprintf("\033[2J\033[H")) // clear screen
		err := cb(ctx)
		if err != nil {
			ctx.Print(fmt.Sprintf(err.Error() + "\n"))
		}
	}
	return nil
}

func (args *Total) printEvaluation(ctx app.Context) error {
	records, err := ctx.RetrieveRecords(args.File...)
	if err != nil {
		return err
	}
	records = service.Query(records, args.toFilter())
	total, _ := func() (Duration, bool) {
		if args.Live {
			return service.HypotheticalTotal(time.Now(), records...)
		}
		return service.Total(records...), false
	}()
	ctx.Print(fmt.Sprintf("Total: %s\n", styler.Duration(total, false)))
	if args.Diff {
		should := service.ShouldTotalSum(records...)
		diff := NewDuration(0, 0).Minus(should).Plus(total)
		ctx.Print(fmt.Sprintf("Should: %s\n", styler.ShouldTotal(should)))
		ctx.Print(fmt.Sprintf("Diff: %s\n", styler.Duration(diff, true)))
	}
	ctx.Print(fmt.Sprintf("(In %d record%s)\n", len(records), func() string {
		if len(records) == 1 {
			return ""
		}
		return "s"
	}()))

	args.WarnArgs.printWarnings(ctx, records)
	return nil
}
