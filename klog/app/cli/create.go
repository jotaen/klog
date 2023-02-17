package cli

import (
	"github.com/jotaen/klog/klog"
	"github.com/jotaen/klog/klog/app"
	"github.com/jotaen/klog/klog/app/cli/lib"
	"github.com/jotaen/klog/klog/parser/reconciling"
)

type Create struct {
	ShouldTotal      klog.ShouldTotal   `name:"should" help:"The should-total of the record"`
	ShouldTotalAlias klog.ShouldTotal   `name:"should-total" hidden:""` // Alias for “canonical” term
	Summary          klog.RecordSummary `name:"summary" short:"s" placeholder:"TEXT" help:"Summary text for the new record"`
	lib.AtDateArgs
	lib.NoStyleArgs
	lib.OutputFileArgs
	lib.WarnArgs
}

func (opt *Create) Help() string {
	return `The new record is inserted into the file at the chronologically correct position.
(Assuming that the records are sorted from oldest to latest.)`
}

func (opt *Create) Run(ctx app.Context) app.Error {
	opt.NoStyleArgs.Apply(&ctx)
	date := opt.AtDate(ctx.Now())
	additionalData := reconciling.AdditionalData{ShouldTotal: opt.GetShouldTotal(), Summary: opt.Summary}
	if additionalData.ShouldTotal == nil {
		ctx.Config().DefaultShouldTotal.Map(func(s klog.ShouldTotal) {
			additionalData.ShouldTotal = s
		})
	}
	return lib.Reconcile(ctx, lib.ReconcileOpts{OutputFileArgs: opt.OutputFileArgs, WarnArgs: opt.WarnArgs},
		[]reconciling.Creator{
			reconciling.NewReconcilerForNewRecord(date, opt.DateFormat(ctx.Config()), additionalData),
		},

		func(reconciler *reconciling.Reconciler) (*reconciling.Result, error) {
			return reconciler.MakeResult()
		},
	)
}

func (opt *Create) GetShouldTotal() klog.ShouldTotal {
	if opt.ShouldTotal != nil {
		return opt.ShouldTotal
	}
	return opt.ShouldTotalAlias
}
