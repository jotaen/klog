package cli

import (
	"github.com/jotaen/klog/klog"
	"github.com/jotaen/klog/klog/app"
	"github.com/jotaen/klog/klog/app/cli/util"
	"github.com/jotaen/klog/klog/parser/reconciling"
)

type Track struct {
	Entry klog.EntrySummary `arg:"" required:"" placeholder:"ENTRY" help:"The new entry to add"`
	util.AtDateArgs
	util.NoStyleArgs
	util.OutputFileArgs
	util.WarnArgs
}

func (opt *Track) Help() string {
	return `The given text is appended to the record as new entry (taken over as is).

Example: klog track '1h work' file.klg

Remember to use 'quotes' if the entry consists of multiple words,
to avoid the text being split or interpreted by your shell.`
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
