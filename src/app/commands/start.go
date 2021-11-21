package commands

import (
	. "github.com/jotaen/klog/src"
	"github.com/jotaen/klog/src/app"
	"github.com/jotaen/klog/src/app/lib"
	"github.com/jotaen/klog/src/parser/reconciling"
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
	return ctx.ReconcileFile(
		opt.OutputFileArgs.File,
		func(reconciler reconciling.Reconciler) (*reconciling.Result, error) {
			return reconciler.AppendEntry(
				func(r Record) bool { return r.Date().IsEqualTo(date) },
				func(r Record) string { return entry },
			)
		},
		func(reconciler reconciling.Reconciler) (*reconciling.Result, error) {
			headline := opt.AtDate(ctx.Now()).ToString()
			lines := []reconciling.InsertableText{
				{headline, 0},
				{entry, 1},
			}
			return reconciler.InsertRecord(date, lines)
		},
	)
}
