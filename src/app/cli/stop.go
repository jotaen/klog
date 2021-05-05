package cli

import (
	. "klog"
	"klog/app"
	"klog/app/cli/lib"
	"klog/parser"
)

type Stop struct {
	lib.AtTimeArgs
	lib.AtDateArgs
	Summary string `name:"summary" short:"s" help:"Text to append to the entry summary"`
	lib.NoStyleArgs
	lib.OutputFileArgs
}

func (opt *Stop) Run(ctx app.Context) error {
	opt.NoStyleArgs.Apply(&ctx)
	date := opt.AtDate(ctx.Now())
	time := opt.AtTime(ctx.Now())
	return lib.ReconcilerChain{
		File: opt.OutputFileArgs.File,
		Ctx:  ctx,
	}.Apply(
		func(pr *parser.ParseResult) (*parser.ReconcileResult, error) {
			reconciler := parser.NewRecordReconciler(pr, func(r Record) bool {
				return r.Date().IsEqualTo(date)
			})
			if reconciler == nil {
				return nil, lib.NotEligibleError{}
			}
			return reconciler.CloseOpenRange(
				func(r Record) (Time, Summary) { return time, Summary(opt.Summary) },
			)
		},
		func(pr *parser.ParseResult) (*parser.ReconcileResult, error) {
			reconciler := parser.NewRecordReconciler(pr, func(r Record) bool {
				return r.Date().IsEqualTo(date.PlusDays(-1))
			})
			if reconciler == nil {
				return nil, lib.NotEligibleError{}
			}
			adjustedTime := func() Time {
				if time.IsTomorrow() {
					return time
				}
				timeTomorrow, _ := time.Add(NewDuration(24, 0))
				return timeTomorrow
			}()
			return reconciler.CloseOpenRange(
				func(r Record) (Time, Summary) { return adjustedTime, Summary(opt.Summary) },
			)
		},
	)
}
