package args

import (
	"github.com/jotaen/klog/klog"
	"github.com/jotaen/klog/klog/app"
)

type SummaryArgs struct {
	SummaryText klog.EntrySummary `name:"summary" short:"s" placeholder:"TEXT" help:"Summary text for the new entry."`
	Resume      bool              `name:"resume" short:"R" help:"Take over summary of last entry (if applicable)."`
	ResumeNth   int               `name:"resume-nth" short:"N" help:"Take over summary of nth entry. If INT is positive, it counts from the start (beginning with '1'); if negative, it counts from the end (beginning with '-1')"`
}

func (args *SummaryArgs) Summary(currentRecord klog.Record, previousRecord klog.Record) (klog.EntrySummary, app.Error) {
	// Check for conflicting flags.
	if args.SummaryText != nil && (args.Resume || args.ResumeNth != 0) {
		return nil, app.NewErrorWithCode(
			app.LOGICAL_ERROR,
			"Conflicting flags: --summary and --resume cannot be used at the same time",
			"",
			nil,
		)
	}
	if args.Resume && args.ResumeNth != 0 {
		return nil, app.NewError(
			"Illegal flag combination",
			"Cannot combine --resume and --resume-nth",
			nil,
		)
	}

	// Return summary flag, if specified.
	if args.SummaryText != nil {
		return args.SummaryText, nil
	}

	// If --resume was specified: return summary of last entry from current record, if
	// it has any entries. Otherwise, return summary of last entry from previous record,
	// if exists.
	if args.Resume {
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
	if args.ResumeNth != 0 {
		if e, ok := findNthEntry(currentRecord, args.ResumeNth); ok {
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
