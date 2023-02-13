package cli

import (
	"github.com/jotaen/klog/klog"
	"github.com/jotaen/klog/klog/app"
	"github.com/jotaen/klog/klog/app/cli/lib"
	"github.com/jotaen/klog/klog/parser/reconciling"
)

type Track struct {
	Entry klog.EntrySummary `arg:"" required:"" placeholder:"ENTRY" help:"The new entry to add"`
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

func (opt *Track) Run(ctx app.Context) app.Error {
	opt.NoStyleArgs.Apply(&ctx)
	now := ctx.Now()
	date, isAutoDate := opt.AtDate(now)
	atDate := reconciling.NewStyled[klog.Date](date, isAutoDate)
	additionalData := reconciling.AdditionalData{}
	ctx.Config().DefaultShouldTotal.Map(func(s klog.ShouldTotal) {
		additionalData.ShouldTotal = s
	})
	return lib.Reconcile(ctx, lib.ReconcileOpts{OutputFileArgs: opt.OutputFileArgs, WarnArgs: opt.WarnArgs},
		[]reconciling.Creator{
			reconciling.NewReconcilerAtRecord(atDate.Value),
			reconciling.NewReconcilerForNewRecord(atDate, additionalData),
		},

		func(reconciler *reconciling.Reconciler) (*reconciling.Result, error) {
			return reconciler.AppendEntry(opt.Entry)
		},
	)
}
