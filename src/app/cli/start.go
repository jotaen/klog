package cli

import (
	. "klog"
	"klog/app"
	"klog/app/cli/lib"
	"klog/parser"
)

type Start struct {
	lib.AtTimeArgs
	lib.AtDateArgs
	lib.NoStyleArgs
	lib.OutputFileArgs
}

func (opt *Start) Run(ctx app.Context) error {
	opt.NoStyleArgs.SetGlobalState()
	date := opt.AtDate(ctx.Now())
	time := opt.AtTime(ctx.Now())
	return applyReconciler(
		opt.OutputFileArgs,
		ctx,
		func(pr *parser.ParseResult) (*parser.ReconcileResult, error) {
			reconciler := parser.NewRecordReconciler(pr, func(r Record) bool {
				return r.Date().IsEqualTo(date) && r.OpenRange() == nil
			})
			if reconciler == nil {
				return nil, app.NewError(
					"No eligible record at date "+date.ToString(),
					"Please make sure the record exists and it doesn’t contain an open-ended time range yet.",
				)
			}
			return reconciler.AppendEntry(
				func(r Record) string { return time.ToString() + " - ?" },
			)
		},
	)
}
