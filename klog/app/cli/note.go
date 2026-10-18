package cli

import (
	"github.com/jotaen/klog/klog"
	"github.com/jotaen/klog/klog/app"
	"github.com/jotaen/klog/klog/app/cli/args"
	"github.com/jotaen/klog/klog/app/cli/helper"
	"github.com/jotaen/klog/klog/parser/reconciling"
)

type Note struct {
	Text       klog.EntrySummary `arg:"" required:"" placeholder:"TEXT" help:"The text to add."`
	Entry      int               `name:"entry" short:"e" placeholder:"INT" help:"Target the nth entry instead of the open time range. If INT is positive, it counts from the start (beginning with '1'); if negative, it counts from the end (beginning with '-1'). If absent, it tries to find and target the open-ended entry."`
	AppendMode string            `name:"append-mode" short:"a" placeholder:"MODE" help:"How to add the text: MODE can be 'inline' (default) or 'newline'." enum:"inline,newline" default:"inline"`
	args.AtDateArgs
	args.NoStyleArgs
	args.WarnArgs
	args.OutputFileArgs
}

func (opt *Note) Help() string {
	return `
By default, it targets the open time range (e.g. '18:00 - ?'). Alternatively, you can choose any entry via '--entry'.

The '--append-mode' flag controls how the text is added:
  - '--append-mode=inline' (default) adds the text to the last line of the summary, delimited by a space character.
  - '--append-mode=newline' adds the text as separate line underneath the summary.

By default, the command targets the record at today’s date, or yesterday’s if there is no record today.
You can otherwise specify a date with '--date'.
`
}

func (opt *Note) Run(ctx app.Context) app.Error {
	opt.NoStyleArgs.Apply(&ctx)
	date := opt.AtDate(ctx.Now())
	creators := []reconciling.Creator{reconciling.NewReconcilerAtRecord(date)}
	if opt.Date == nil {
		// If no date had been specified explicitly, try to fall back to yesterday’s record.
		creators = append(creators, reconciling.NewReconcilerAtRecord(date.PlusDays(-1)))
	}
	summary := opt.Text
	if opt.AppendMode == "newline" && summary[0] != "" {
		summary = append(klog.EntrySummary{""}, summary...)
	}
	return helper.Reconcile(ctx, helper.ReconcileOpts{OutputFileArgs: opt.OutputFileArgs, WarnArgs: opt.WarnArgs},
		creators,
		func(reconciler *reconciling.Reconciler) error {
			if opt.Entry != 0 {
				return reconciler.AppendToNthEntry(opt.Entry, summary)
			}
			return reconciler.AppendToOpenRange(summary)
		},
	)
}
