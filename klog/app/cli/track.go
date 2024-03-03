package cli

import (
	"github.com/jotaen/klog/klog"
	"github.com/jotaen/klog/klog/app"
	"github.com/jotaen/klog/klog/app/cli/util"
	"github.com/jotaen/klog/klog/parser/reconciling"
)

type Track struct {
	Entry klog.EntrySummary `arg:"" required:"" placeholder:"ENTRY" help:"The new entry to add."`
	util.AtDateArgs
	util.NoStyleArgs
	util.WarnArgs
	util.OutputFileArgs
}

func (opt *Track) Help() string {
	return `
The given text is appended to the record as new entry (taken over as is, i.e. including the entry summary). Example invocations:

    klog track '1h' file.klg
    klog track '15:00 - 16:00 Went out running' file.klg
    klog track '6h30m #work' file.klg

It uses the record at today’s date for the new entry, or creates a new record if there no record at today’s date.
You can otherwise specify a date with '--date'.

Remember to use 'quotes' if the entry consists of multiple words, to avoid the text being split or otherwise pre-processed by your shell.
There is still one quirk: if you want to track a negative duration, you have to escape the leading minus with a backslash, e.g. '\-45m lunch break', to prevent it from being mistakenly interpreted as a flag.
`
}

func (opt *Track) Run(ctx app.Context) app.Error {
	opt.NoStyleArgs.Apply(&ctx)
	now := ctx.Now()
	date := opt.AtDate(now)
	additionalData := reconciling.AdditionalData{}
	ctx.Config().DefaultShouldTotal.Unwrap(func(s klog.ShouldTotal) {
		additionalData.ShouldTotal = s
	})
	return util.Reconcile(ctx, util.ReconcileOpts{OutputFileArgs: opt.OutputFileArgs, WarnArgs: opt.WarnArgs},
		[]reconciling.Creator{
			reconciling.NewReconcilerAtRecord(date),
			reconciling.NewReconcilerForNewRecord(date, opt.DateFormat(ctx.Config()), additionalData),
		},

		func(reconciler *reconciling.Reconciler) error {
			return reconciler.AppendEntry(opt.Entry)
		},
	)
}
