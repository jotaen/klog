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
	lib.NoStyleArgs
	lib.OutputFileArgs
}

func (opt *Stop) Run(ctx app.Context) error {
	opt.NoStyleArgs.SetGlobalState()
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
				return nil, app.NewError(
					"No eligible record at date "+date.ToString(),
					"Please make sure the record exists and it contains an open-ended time range "+
						"which start time is prior to your desired end time.",
					nil,
				)
			}
			return reconciler.CloseOpenRange(
				func(r Record) Time { return time },
			)
		},
	)
}
