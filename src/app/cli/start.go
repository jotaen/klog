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
	return reconcile(
		opt.OutputFileArgs,
		ctx,
		func(pr *parser.ParseResult) (*parser.Reconciler, error) {
			return parser.NewRecordReconciler(pr,
				errors.New("No eligible record at date "+date.ToString()),
				func(r Record) bool {
					return r.Date().IsEqualTo(date) &&
						r.OpenRange() == nil
				})
		},
		func(r *parser.Reconciler) (Record, string, error) {
			return r.AppendEntry(
				func(r Record) string { return time.ToString() + " - ?" },
			)
		},
	)
}
