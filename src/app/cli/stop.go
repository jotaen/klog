package cli

import (
	"errors"
	. "klog"
	"klog/app"
	"klog/app/cli/lib"
	"klog/parser"
)

type Stop struct {
	lib.AtDateArgs
	lib.OutputFileArgs
}

func (opt *Stop) Run(ctx app.Context) error {
	date := opt.AtDate(ctx.Now())
	time := NewTimeFromTime(ctx.Now())
	return reconcile(
		opt.OutputFileArgs,
		ctx,
		errors.New("No record (with open time range) at date "+date.ToString()),
		func(r Record) bool { return r.Date().IsEqualTo(date) && r.OpenRange() != nil },
		func(r *parser.Reconciler) (Record, string, error) {
			return r.CloseOpenRange(
				func(r Record) Time { return time },
			)
		},
	)
}
