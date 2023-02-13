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

func (opt *Start) Run(ctx app.Context) app.Error {
	opt.NoStyleArgs.Apply(&ctx)
	now := ctx.Now()
	date, isAutoDate := opt.AtDate(now)
	time, isAutoTime, err := opt.AtTime(now, ctx.Config())
	if err != nil {
		return err
	}
	atDate := reconciling.NewStyled[klog.Date](date, isAutoDate)
	startTime := reconciling.NewStyled[klog.Time](time, isAutoTime)
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
			return reconciler.StartOpenRange(startTime, opt.Summary)
		},
	)
}
