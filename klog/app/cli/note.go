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
	Entry      int               `name:"entry" short:"e" placeholder:"INT" help:"Target the nth entry instead of the open time range. If INT is positive, it counts from the start (beginning with '1'); if negative, it counts from the end (beginning with '-1')."`
	AppendMode string            `name:"append-mode" short:"a" placeholder:"MODE" help:"How to add the text. MODE can be 'inline' (default) or 'newline'." enum:"inline,newline" default:"inline"`
	args.AtDateArgs
	args.NoStyleArgs
	args.WarnArgs
	args.OutputFileArgs
}

func (opt *Note) Help() string {
	return `
Adds text to the summary of an entry.
By default, it targets the open time range (e.g. '18:00 - ?'). You can target any other entry via '--entry'.

With '--append-mode=inline' (the default), the text is appended to the last line of the summary.
With '--append-mode=newline', it is added as a new line underneath the summary.
Tags in the text become tags of the entry.

By default, it targets the record at today’s date, or yesterday’s if there is no record today.
You can otherwise specify a date with '--date'.
`
}

func (opt *Note) Run(ctx app.Context) app.Error {
	opt.NoStyleArgs.Apply(&ctx)
	date := opt.AtDate(ctx.Now())
	creators := []reconciling.Creator{reconciling.NewReconcilerAtRecord(date)}
	if opt.Date == nil {
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
