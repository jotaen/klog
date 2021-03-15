package cli

import (
	"errors"
	. "klog"
	"klog/app"
	"klog/app/cli/lib"
	"klog/parser"
)

type Stop struct {
	lib.AtTimeArgs
	lib.AtDateArgs
	lib.OutputFileArgs
}

func (opt *Stop) Run(ctx app.Context) error {
	date := opt.AtDate(ctx.Now())
	time := opt.AtTime(ctx.Now())
	return applyReconciler(
		opt.OutputFileArgs,
		ctx,
		func(pr *parser.ParseResult) (*parser.ReconcileResult, error) {
			reconciler := parser.NewRecordReconciler(pr, func(r Record) bool {
				return r.Date().IsEqualTo(date) &&
					r.OpenRange() != nil &&
					time.IsAfterOrEqual(r.OpenRange().Start())
			})
			if reconciler == nil {
				return nil, errors.New("No eligible record at date " + date.ToString())
			}
			return reconciler.CloseOpenRange(
				func(r Record) Time { return time },
			)
		},
	)
}
