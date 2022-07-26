package cli

import (
	. "github.com/jotaen/klog/src"
	"github.com/jotaen/klog/src/app"
	"github.com/jotaen/klog/src/app/cli/lib"
	"github.com/jotaen/klog/src/parser"
	"github.com/jotaen/klog/src/parser/reconciling"
)

type Track struct {
	Entry EntrySummary `arg:"" required:"" placeholder:"ENTRY" help:"The new entry to add"`
	lib.AtDateArgs
	lib.NoStyleArgs
	lib.OutputFileArgs
	lib.WarnArgs
}

func (opt *Track) Help() string {
	return `The text of the new entry is taken over as is and appended to the record.

Example: klog track '1h work' file.klg

Remember to use 'quotes' if the entry consists of multiple words,
and to avoid the text being processed by your shell.`
}

func (opt *Track) Run(ctx app.Context) error {
	opt.NoStyleArgs.Apply(&ctx)
	now := ctx.Now()
	date, _ := opt.AtDate(now)
	return lib.Reconcile(ctx, lib.ReconcileOpts{OutputFileArgs: opt.OutputFileArgs, WarnArgs: opt.WarnArgs},
		[]reconciling.Creator{
			func(parsedRecords []parser.ParsedRecord) *reconciling.Reconciler {
				return reconciling.NewReconcilerAtRecord(parsedRecords, date)
			},
			func(parsedRecords []parser.ParsedRecord) *reconciling.Reconciler {
				return reconciling.NewReconcilerForNewRecord(parsedRecords, reconciling.RecordParams{Date: date})
			},
		},

		func(reconciler *reconciling.Reconciler) (*reconciling.Result, error) {
			return reconciler.AppendEntry(opt.Entry)
		},
	)
}
