package cli

import (
	"errors"
	. "klog"
	"klog/app"
	"klog/app/cli/lib"
	"klog/parser"
)

type Start struct {
	lib.AtDateArgs
	lib.OutputFileArgs
}

func (opt *Start) Run(ctx app.Context) error {
	date := opt.AtDate(ctx.Now())
	time := NewTimeFromTime(ctx.Now())
	return reconcile(
		opt.OutputFileArgs,
		ctx,
		errors.New("No record at date "+date.ToString()),
		func(r Record) bool { return r.Date().IsEqualTo(date) },
		func(r *parser.Reconciler) (Record, string, error) {
			return r.AppendEntry(
				func(r Record) string { return time.ToString() + "-?" },
			)
		},
	)
}
