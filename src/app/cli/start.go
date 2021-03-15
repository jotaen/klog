package cli

import (
	"errors"
	. "klog"
	"klog/app"
	"klog/app/cli/lib"
	"klog/parser"
)

type Start struct {
	lib.AtTimeArgs
	lib.AtDateArgs
	lib.OutputFileArgs
}

func (opt *Start) Run(ctx app.Context) error {
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
				return nil, errors.New("No eligible record at date " + date.ToString())
			}
			return reconciler.AppendEntry(
				func(r Record) string { return time.ToString() + " - ?" },
			)
		},
	)
}
