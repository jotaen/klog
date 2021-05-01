package cli

import (
	"errors"
	. "klog"
	"klog/app"
	"klog/app/cli/lib"
	"klog/parser"
	"strings"
)

type Track struct {
	lib.AtDateArgs
	Entry string `arg required help:"The new entry to add, which may optionally contain summary text. Remember to 'quote' to avoid shell processing."`
	lib.NoStyleArgs
	lib.OutputFileArgs
}

func (opt *Track) Run(ctx app.Context) error {
	opt.NoStyleArgs.SetGlobalState()
	date := opt.AtDate(ctx.Now())
	value := sanitiseQuotedLeadingDash(opt.Entry)
	return ReconcilerApplicator{
		file: opt.OutputFileArgs.File,
		ctx:  ctx,
	}.apply(
		func(pr *parser.ParseResult) (*parser.ReconcileResult, error) {
			reconciler := parser.NewRecordReconciler(pr, func(r Record) bool {
				return r.Date().IsEqualTo(date)
			})
			if reconciler == nil {
				return nil, errors.New("No record at date " + date.ToString())
			}
			return reconciler.AppendEntry(func(r Record) string { return value })
		},
	)
}

func sanitiseQuotedLeadingDash(text string) string {
	// When passing entries like `-45m` the leading dash must be escaped
	// otherwise it’s treated like a flag. Therefore we have to remove
	// the potential escaping backslash.
	return strings.TrimPrefix(text, "\\")
}

type ReconcilerApplicator struct {
	file string
	ctx  app.Context
}

func (a ReconcilerApplicator) apply(
	reconcile func(*parser.ParseResult) (*parser.ReconcileResult, error),
) error {
	pr, err := a.ctx.ReadFileInput(a.file)
	if err != nil {
		return err
	}
	result, err := reconcile(pr)
	if err != nil {
		return err
	}
	err = a.ctx.WriteFile(a.file, result.NewText)
	if err != nil {
		return err
	}
	a.ctx.Print("\n" + parser.SerialiseRecords(&lib.Styler, result.NewRecord) + "\n")
	return nil
}
