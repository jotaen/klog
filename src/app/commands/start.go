package commands

import (
	. "github.com/jotaen/klog/src"
	"github.com/jotaen/klog/src/app"
	"github.com/jotaen/klog/src/app/lib"
	"github.com/jotaen/klog/src/parser/lineparsing"
	"github.com/jotaen/klog/src/parser/reconciler"
)

type Start struct {
	lib.AtTimeArgs
	lib.AtDateArgs
	Summary string `name:"summary" short:"s" help:"Summary text for this entry"`
	lib.NoStyleArgs
	lib.OutputFileArgs
}

func (opt *Start) Help() string {
	return `A new open-ended entry is appended to the record, e.g. 14:00-?.
The start time is the current time (or whatever is specified by --time).`
}

func (opt *Start) Run(ctx app.Context) error {
	opt.NoStyleArgs.Apply(&ctx)
	date := opt.AtDate(ctx.Now())
	time := opt.AtTime(ctx.Now())
	entry := func() string {
		summary := ""
		if opt.Summary != "" {
			summary += " " + opt.Summary
		}
		return time.ToString() + " - ?" + summary
	}()
	return ctx.ReconcileFile(opt.OutputFileArgs.File,
		func(records []Record, blocks []lineparsing.Block) (*reconciler.ReconcileResult, error) {
			entryReconciler := reconciler.NewEntryReconciler(records, blocks, func(r Record) bool {
				return r.Date().IsEqualTo(date)
			})
			if entryReconciler == nil {
				return nil, reconciler.NotEligibleError{}
			}
			return entryReconciler.AppendEntry(func(r Record) string {
				return entry
			})
		},
		func(records []Record, blocks []lineparsing.Block) (*reconciler.ReconcileResult, error) {
			recordReconciler := reconciler.NewRecordReconciler(records, blocks, date)
			headline := opt.AtDate(ctx.Now()).ToString()
			lines := []reconciler.InsertableText{
				{headline, 0},
				{entry, 1},
			}
			return recordReconciler.InsertBlock(lines)
		},
	)
}
