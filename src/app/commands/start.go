package commands

import (
	. "github.com/jotaen/klog/src"
	"github.com/jotaen/klog/src/app"
	"github.com/jotaen/klog/src/app/lib"
	"github.com/jotaen/klog/src/parser"
	"github.com/jotaen/klog/src/parser/parsing"
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
	return lib.ReconcilerChain{
		File: opt.OutputFileArgs.File,
		Ctx:  ctx,
	}.Apply(
		func(pr *parser.ParseResult) (*parser.ReconcileResult, error) {
			reconciler := parser.NewRecordReconciler(pr, func(r Record) bool {
				return r.Date().IsEqualTo(date)
			})
			if reconciler == nil {
				return nil, lib.NotEligibleError{}
			}
			return reconciler.AppendEntry(func(r Record) string {
				return entry
			})
		},
		func(pr *parser.ParseResult) (*parser.ReconcileResult, error) {
			reconciler := parser.NewBlockReconciler(pr, date)
			headline := opt.AtDate(ctx.Now()).ToString()
			lines := []parsing.Text{
				{headline, 0},
				{entry, 1},
			}
			return reconciler.InsertBlock(lines)
		},
	)
}
