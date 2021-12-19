package commands

import (
	"github.com/jotaen/klog/src/app"
	"github.com/jotaen/klog/src/app/lib"
	"github.com/jotaen/klog/src/parser"
	"github.com/jotaen/klog/src/parser/reconciling"
)

type Start struct {
	lib.AtDateAndTimeArgs
	Summary string `name:"summary" short:"s" help:"Summary text for this entry"`
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
			func(parsedRecords []parser.ParsedRecord) *reconciling.Reconciler {
				return reconciling.NewReconcilerAtRecord(parsedRecords, date)
			},
			func(parsedRecords []parser.ParsedRecord) *reconciling.Reconciler {
				return reconciling.NewReconcilerAtNewRecord(parsedRecords, date, nil)
			},
		},

		func(reconciler *reconciling.Reconciler) (*reconciling.Result, error) {
			return reconciler.StartOpenRange(time, opt.Summary)
		},
	)
}
