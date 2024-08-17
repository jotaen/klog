package cli

import (
	"github.com/jotaen/klog/klog"
	"github.com/jotaen/klog/klog/app"
	"github.com/jotaen/klog/klog/app/cli/util"
	"github.com/jotaen/klog/klog/parser/reconciling"
	"github.com/jotaen/klog/klog/parser/txt"
	"github.com/jotaen/klog/klog/service"
)

type Start struct {
	util.SummaryArgs
	util.AtDateAndTimeArgs
	util.NoStyleArgs
	util.WarnArgs
	util.OutputFileArgs
}

func (opt *Start) Help() string {
	return `
This appends a new open-ended entry to the record.

By default, it uses the record at today’s date for the new entry, or creates a new record if there no record at today’s date.
You can otherwise specify a date with '--date'.

Unless the '--time' flag is specified, it defaults to the current time as start time.
If you prefer your time to be rounded, you can use the '--round' flag.

You can either assign a summary text for the new entry via the '--summary' flag, or you can use the '--resume' flag to automatically take over the entry summary of the last entry.
Note that '--resume' will fall back to the last record, if the current record doesn’t contain any entries yet.
`
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
	ctx.Config().DefaultShouldTotal.Unwrap(func(s klog.ShouldTotal) {
		additionalData.ShouldTotal = s
	})

	spy := PreviousRecordSpy{}
	return util.Reconcile(ctx, util.ReconcileOpts{OutputFileArgs: opt.OutputFileArgs, WarnArgs: opt.WarnArgs},
		[]reconciling.Creator{
			spy.phonyCreator(date),
			reconciling.NewReconcilerAtRecord(date),
			reconciling.NewReconcilerForNewRecord(date, opt.DateFormat(ctx.Config()), additionalData),
		},

		func(reconciler *reconciling.Reconciler) error {
			summary, sErr := opt.Summary(reconciler.Record, spy.PreviousRecord)
			if sErr != nil {
				return sErr
			}
			return reconciler.StartOpenRange(time, opt.TimeFormat(ctx.Config()), summary)
		},
	)
}

type PreviousRecordSpy struct {
	PreviousRecord klog.Record
}

// phonyCreator is a no-op “pass-through” creator, whose only purpose it is to hook into
// the reconciler-creation mechanism, to get a handle on the records for determining
// the previous record.
func (p *PreviousRecordSpy) phonyCreator(currentDate klog.Date) reconciling.Creator {
	return func(records []klog.Record, _ []txt.Block) *reconciling.Reconciler {
		for _, r := range service.Sort(records, false) {
			if r.Date().IsAfterOrEqual(currentDate) {
				continue
			}
			p.PreviousRecord = r
			return nil
		}
		return nil
	}
}
