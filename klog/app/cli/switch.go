package cli

import (
	"github.com/jotaen/klog/klog/app"
	"github.com/jotaen/klog/klog/app/cli/util"
	"github.com/jotaen/klog/klog/parser/reconciling"
)

type Switch struct {
	util.SummaryArgs
	util.AtDateAndTimeArgs
	util.NoStyleArgs
	util.WarnArgs
	util.OutputFileArgs
}

func (opt *Switch) Help() string {
	return `
Closes a previously ongoing activity (i.e., open time range), and starts a new one.
This is basically a convenience for doing 'klog stop' and 'klog start' – however, in contrast to issuing both commands separately, 'klog switch' guarantees that the end time of the previous activity will be the same as the start time for the new entry.

By default, it uses the record at today’s date for the new entry. You can otherwise specify a date with '--date'.

Unless the '--time' flag is specified, it defaults to the current time as start/stop time.
If you prefer your time to be rounded, you can use the '--round' flag.
`
}

func (opt *Switch) Run(ctx app.Context) app.Error {
	opt.NoStyleArgs.Apply(&ctx)
	now := ctx.Now()
	date := opt.AtDate(now)
	time, tErr := opt.AtTime(now, ctx.Config())
	if tErr != nil {
		return tErr
	}

	return util.Reconcile(ctx, util.ReconcileOpts{OutputFileArgs: opt.OutputFileArgs, WarnArgs: opt.WarnArgs},
		[]reconciling.Creator{
			reconciling.NewReconcilerAtRecord(date),
		},

		func(reconciler *reconciling.Reconciler) error {
			return reconciler.CloseOpenRange(time, opt.TimeFormat(ctx.Config()), nil)
		},
		func(reconciler *reconciling.Reconciler) error {
			summary, sErr := opt.Summary(reconciler.Record, nil)
			if sErr != nil {
				return sErr
			}
			return reconciler.StartOpenRange(time, opt.TimeFormat(ctx.Config()), summary)
		},
	)
}
