package cli

import (
	"github.com/jotaen/klog/klog"
	"github.com/jotaen/klog/klog/app"
	"github.com/jotaen/klog/klog/app/cli/lib"
	"github.com/jotaen/klog/klog/parser/reconciling"
)

type Start struct {
	Summary klog.EntrySummary `name:"summary" short:"s" placeholder:"TEXT" help:"Summary text for this entry"`
	lib.AtDateAndTimeArgs
	lib.NoStyleArgs
	lib.OutputFileArgs
	lib.WarnArgs
}

func (opt *Start) Help() string {
	return `A new open-ended entry is appended to the record, e.g. 14:00-?.
The start time is the current time (or whatever is specified by --time).`
}

func (opt *Start) Run(ctx app.Context) error {
	opt.NoStyleArgs.Apply(&ctx)
	now := ctx.Now()
	date, _ := opt.AtDate(now)
	time, _, err := opt.AtTime(now)
	if err != nil {
		return err
	}
	return lib.Reconcile(ctx, lib.ReconcileOpts{OutputFileArgs: opt.OutputFileArgs, WarnArgs: opt.WarnArgs},
		[]reconciling.Creator{
			reconciling.NewReconcilerAtRecord(date),
			reconciling.NewReconcilerForNewRecord(reconciling.RecordParams{Date: date}),
		},

		func(reconciler *reconciling.Reconciler) (*reconciling.Result, error) {
			return reconciler.StartOpenRange(time, opt.Summary)
		},
	)
}
