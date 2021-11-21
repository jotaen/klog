package commands

import (
	. "github.com/jotaen/klog/src"
	"github.com/jotaen/klog/src/app"
	"github.com/jotaen/klog/src/app/lib"
	"github.com/jotaen/klog/src/parser/reconciling"
	"strings"
)

type Track struct {
	lib.AtDateArgs
	Entry string `arg required help:"The new entry to add"`
	lib.NoStyleArgs
	lib.OutputFileArgs
}

func (opt *Track) Help() string {
	return `The text of the new entry is taken over as is and appended to the record.

Example: klog track '1h work' file.klg

Remember to use 'quotes' if the entry consists of multiple words,
and to avoid the text being processed by your shell.`
}

func (opt *Track) Run(ctx app.Context) error {
	opt.NoStyleArgs.Apply(&ctx)
	date := opt.AtDate(ctx.Now())
	value := sanitiseQuotedLeadingDash(opt.Entry)
	return ctx.ReconcileFile(
		opt.OutputFileArgs.File,
		func(reconciler reconciling.Reconciler) (*reconciling.Result, error) {
			return reconciler.AppendEntry(
				func(r Record) bool { return r.Date().IsEqualTo(date) },
				func(r Record) string { return value })
		},
		func(reconciler reconciling.Reconciler) (*reconciling.Result, error) {
			headline := opt.AtDate(ctx.Now()).ToString()
			lines := []reconciling.InsertableText{
				{headline, 0},
				{value, 1},
			}
			return reconciler.InsertRecord(date, lines)
		},
	)
}

func sanitiseQuotedLeadingDash(text string) string {
	// When passing entries like `-45m` the leading dash must be escaped
	// otherwise it’s treated like a flag. Therefore we have to remove
	// the potential escaping backslash.
	return strings.TrimPrefix(text, "\\")
}
