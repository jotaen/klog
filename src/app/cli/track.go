package cli

import (
	. "klog"
	"klog/app"
	"klog/app/cli/lib"
	"klog/parser"
	"klog/parser/parsing"
	"strings"
)

type Track struct {
	lib.AtDateArgs
	Entry string `arg required help:"The new entry to add, which may optionally contain summary text. Remember to 'quote' to avoid shell processing."`
	lib.NoStyleArgs
	lib.OutputFileArgs
}

func (opt *Track) Run(ctx app.Context) error {
	opt.NoStyleArgs.Apply(&ctx)
	date := opt.AtDate(ctx.Now())
	value := sanitiseQuotedLeadingDash(opt.Entry)
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
			return reconciler.AppendEntry(func(r Record) string { return value })
		},
		func(pr *parser.ParseResult) (*parser.ReconcileResult, error) {
			reconciler := parser.NewBlockReconciler(pr, date)
			headline := opt.AtDate(ctx.Now()).ToString()
			lines := []parsing.Text{
				{headline, 0},
				{value, 1},
			}
			return reconciler.InsertBlock(lines)
		},
	)
}

func sanitiseQuotedLeadingDash(text string) string {
	// When passing entries like `-45m` the leading dash must be escaped
	// otherwise it’s treated like a flag. Therefore we have to remove
	// the potential escaping backslash.
	return strings.TrimPrefix(text, "\\")
}
