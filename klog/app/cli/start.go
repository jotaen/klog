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
	SummaryText klog.EntrySummary `name:"summary" short:"s" placeholder:"TEXT" help:"Summary text for this entry."`
	Resume      bool              `name:"resume" short:"R" help:"Take over summary of last entry (if applicable). If the target record is new or empty, it looks at the previous record."`
	ResumeNth   int               `name:"resume-nth" short:"N" help:"Take over summary of nth entry. If INT is positive, it counts from the start (beginning with '1'); if negative, it counts from the end (beginning with '-1')"`
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

func (opt *Start) Summary(currentRecord klog.Record, previousRecord klog.Record) (klog.EntrySummary, app.Error) {
	// Check for conflicting flags.
	if opt.SummaryText != nil && (opt.Resume || opt.ResumeNth != 0) {
		return nil, app.NewErrorWithCode(
			app.LOGICAL_ERROR,
			"Conflicting flags: --summary and --resume cannot be used at the same time",
			"",
			nil,
		)
	}
	if opt.Resume && opt.ResumeNth != 0 {
		return nil, app.NewError(
			"Illegal flag combination",
			"Cannot combine --resume and --resume-nth",
			nil,
		)
	}

	// Return summary flag, if specified.
	if opt.SummaryText != nil {
		return opt.SummaryText, nil
	}

	// If --resume was specified: return summary of last entry from current record, if
	// it has any entries. Otherwise, return summary of last entry from previous record,
	// if exists.
	if opt.Resume {
		if e, ok := findNthEntry(currentRecord, -1); ok {
			return e.Summary(), nil
		}
		if previousRecord != nil {
			if e, ok := findNthEntry(previousRecord, -1); ok {
				return e.Summary(), nil
			}
		}
		return nil, nil
	}

	// If --resume-nth was specified: return summary of nth-entry. In contrast to --resume,
	// don’t fall back to previous record, as that would be unintuitive here.
	if opt.ResumeNth != 0 {
		if e, ok := findNthEntry(currentRecord, opt.ResumeNth); ok {
			return e.Summary(), nil
		}
		return nil, app.NewError(
			"No such entry",
			"",
			nil,
		)
	}

	return nil, nil
}

func findNthEntry(r klog.Record, nr int) (klog.Entry, bool) {
	entriesCount := len(r.Entries())
	i := func() int {
		if nr > 0 {
			return nr - 1
		}
		return entriesCount + nr
	}()
	if i < 0 || i > entriesCount-1 {
		return klog.Entry{}, false
	}
	return r.Entries()[i], true
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
