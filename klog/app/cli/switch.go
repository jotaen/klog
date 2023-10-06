package cli

import (
	"github.com/jotaen/klog/klog"
	"github.com/jotaen/klog/klog/app"
	"github.com/jotaen/klog/klog/app/cli/lib"
	"github.com/jotaen/klog/klog/parser/reconciling"
)

type Switch struct {
	SummaryText klog.EntrySummary `name:"summary" short:"s" placeholder:"TEXT" help:"Summary text for the new entry"`
	lib.AtDateAndTimeArgs
	lib.NoStyleArgs
	lib.OutputFileArgs
	lib.WarnArgs
}

func (opt *Switch) Help() string {
	return `Closes a previously ongoing activity (i.e., open time range), and starts a new one.

The end time of the previous activity will be the same as the start time for the new entry.
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

	return lib.Reconcile(ctx, lib.ReconcileOpts{OutputFileArgs: opt.OutputFileArgs, WarnArgs: opt.WarnArgs},
		[]reconciling.Creator{
			reconciling.NewReconcilerAtRecord(date),
		},

		func(reconciler *reconciling.Reconciler) error {
			return reconciler.CloseOpenRange(time, opt.TimeFormat(ctx.Config()), nil)
		},
		func(reconciler *reconciling.Reconciler) error {
			return reconciler.StartOpenRange(time, opt.TimeFormat(ctx.Config()), opt.SummaryText)
		},
	)
}
