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
	Summary string `name:"summary" short:"s" help:"Summary text for this entry"`
	lib.NoStyleArgs
	lib.OutputFileArgs
}

func (opt *Start) Run(ctx app.Context) error {
	opt.NoStyleArgs.SetGlobalState()
	date := opt.AtDate(ctx.Now())
	time := opt.AtTime(ctx.Now())
	return ReconcilerApplicator{
		file: opt.OutputFileArgs.File,
		ctx:  ctx,
	}.apply(
		func(pr *parser.ParseResult) (*parser.ReconcileResult, error) {
			reconciler := parser.NewRecordReconciler(pr, func(r Record) bool {
				return r.Date().IsEqualTo(date) && r.OpenRange() == nil
			})
			if reconciler == nil {
				return nil, app.NewError(
					"No eligible record at date "+date.ToString(),
					"Please make sure the record exists and it doesn’t contain an open-ended time range yet.",
					nil,
				)
			}
			return reconciler.AppendEntry(
				func(r Record) string {
					summary := ""
					if opt.Summary != "" {
						summary += " " + opt.Summary
					}
					return time.ToString() + " - ?" + summary
				},
			)
		},
	)
}
