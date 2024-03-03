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
	SummaryText klog.EntrySummary `name:"summary" short:"s" placeholder:"TEXT" help:"Summary text for this entry"`
	Resume      bool              `name:"resume" short:"R" help:"Take over summary of last entry (if applicable)"`
	util.AtDateAndTimeArgs
	util.NoStyleArgs
	util.OutputFileArgs
	util.WarnArgs
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

func (opt *Start) Summary(currentRecord klog.Record, previousRecord klog.Record) (klog.EntrySummary, app.Error) {
	// Check for conflicting flags.
	if opt.SummaryText != nil && opt.Resume {
		return nil, app.NewErrorWithCode(
			app.LOGICAL_ERROR,
			"Conflicting flags: --summary and --resume cannot be used at the same time",
			"",
			nil,
		)
	}

	// Return summary flag, if specified.
	if opt.SummaryText != nil {
		return opt.SummaryText, nil
	}

	// Skip if resume flag wasn’t specified.
	if !opt.Resume {
		return nil, nil
	}

	// Return summary of last entry from current record, if it has any entries.
	if len(currentRecord.Entries()) > 0 {
		return lastEntrySummary(currentRecord), nil
	}

	// Return summary of last entry from previous record, if exists.
	if previousRecord != nil {
		return lastEntrySummary(previousRecord), nil
	}

	return nil, nil
}

func lastEntrySummary(r klog.Record) klog.EntrySummary {
	entriesCount := len(r.Entries())
	if entriesCount == 0 {
		return nil
	}
	return r.Entries()[entriesCount-1].Summary()
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
