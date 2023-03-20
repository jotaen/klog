package cli

import (
	"github.com/jotaen/klog/klog"
	"github.com/jotaen/klog/klog/app"
	"github.com/jotaen/klog/klog/app/cli/lib"
	"github.com/jotaen/klog/klog/parser/reconciling"
)

type Start struct {
	SummaryText klog.EntrySummary `name:"summary" short:"s" placeholder:"TEXT" help:"Summary text for this entry"`
	Resume      bool              `name:"resume" short:"R" help:"Take over summary of last entry (if applicable)"`
	lib.AtDateAndTimeArgs
	lib.NoStyleArgs
	lib.OutputFileArgs
	lib.WarnArgs
}

func (opt *Start) Help() string {
	return `A new open-ended entry is appended to the record, e.g. 14:00-?.

If the --time flag is not specified, it defaults to the current time as start time. In the latter case, the time can be rounded via --round.`
}

func (opt *Start) Run(ctx app.Context) app.Error {
	opt.NoStyleArgs.Apply(&ctx)
	now := ctx.Now()
	date := opt.AtDate(now)
	time, tErr := opt.AtTime(now, ctx.Config())
	if tErr != nil {
		return tErr
	}
	additionalData := reconciling.AdditionalData{}
	ctx.Config().DefaultShouldTotal.Map(func(s klog.ShouldTotal) {
		additionalData.ShouldTotal = s
	})
	return lib.Reconcile(ctx, lib.ReconcileOpts{OutputFileArgs: opt.OutputFileArgs, WarnArgs: opt.WarnArgs},
		[]reconciling.Creator{
			reconciling.NewReconcilerAtRecord(date),
			reconciling.NewReconcilerForNewRecord(date, opt.DateFormat(ctx.Config()), additionalData),
		},

		func(reconciler *reconciling.Reconciler) (*reconciling.Result, error) {
			summary, sErr := opt.Summary(reconciler.Record)
			if sErr != nil {
				return nil, sErr
			}
			return reconciler.StartOpenRange(time, opt.TimeFormat(ctx.Config()), summary)
		},
	)
}

func (opt *Start) Summary(r klog.Record) (klog.EntrySummary, app.Error) {
	if opt.SummaryText != nil && opt.Resume {
		return nil, app.NewErrorWithCode(
			app.LOGICAL_ERROR,
			"Conflicting flags: --summary and --resume cannot be used at the same time",
			"",
			nil,
		)
	}
	if opt.SummaryText != nil {
		return opt.SummaryText, nil
	}
	entriesCount := len(r.Entries())
	if opt.Resume && entriesCount > 0 {
		return r.Entries()[entriesCount-1].Summary(), nil
	}
	return nil, nil
}
