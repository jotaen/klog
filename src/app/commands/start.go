package commands

import (
	"github.com/jotaen/klog/src/app"
	"github.com/jotaen/klog/src/app/lib"
	"github.com/jotaen/klog/src/parser"
	"github.com/jotaen/klog/src/parser/reconciling"
)

type Start struct {
	lib.AtTimeArgs
	lib.AtDateArgs
	Summary string `name:"summary" short:"s" help:"Summary text for this entry"`
	lib.NoStyleArgs
	lib.OutputFileArgs
}

func (opt *Start) Help() string {
	return `A new open-ended entry is appended to the record, e.g. 14:00-?.
The start time is the current time (or whatever is specified by --time).`
}

func (opt *Start) Run(ctx app.Context) error {
	opt.NoStyleArgs.Apply(&ctx)
	date := opt.AtDate(ctx.Now())
	time := opt.AtTime(ctx.Now())
	return ctx.ReconcileFile(
		opt.OutputFileArgs.File,

		[]reconciling.Creator{
			func(parsedRecords []parser.ParsedRecord) *reconciling.Reconciler {
				return reconciling.NewReconcilerAtRecord(parsedRecords, date)
			},
			func(parsedRecords []parser.ParsedRecord) *reconciling.Reconciler {
				return reconciling.NewReconcilerAtNewRecord(parsedRecords, date, nil)
			},
		},

		// Start new open range at that date.
		func(reconciler *reconciling.Reconciler) (*reconciling.Result, error) {
			return reconciler.StartOpenRange(date, time, opt.Summary)
		},
	)
}
