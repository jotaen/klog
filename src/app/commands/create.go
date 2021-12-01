package commands

import (
	. "github.com/jotaen/klog/src"
	"github.com/jotaen/klog/src/app"
	"github.com/jotaen/klog/src/app/lib"
	"github.com/jotaen/klog/src/parser/reconciling"
)

type Create struct {
	ShouldTotal Duration `name:"should" help:"The should-total of the record"`
	lib.AtDateArgs
	lib.NoStyleArgs
	lib.OutputFileArgs
}

func (opt *Create) Help() string {
	return `The new record is inserted into the file at the chronologically correct position.
(Assuming that the records are sorted from oldest to latest.)`
}

func (opt *Create) Run(ctx app.Context) error {
	opt.NoStyleArgs.Apply(&ctx)
	date := opt.AtDate(ctx.Now())
	lines, err := func() ([]reconciling.InsertableText, error) {
		headline := opt.AtDate(ctx.Now()).ToString()
		if opt.ShouldTotal != nil {
			headline += " (" + opt.ShouldTotal.ToString() + "!)"
		}
		return []reconciling.InsertableText{
			{Text: headline, Indentation: 0},
		}, nil
	}()
	if err != nil {
		return err
	}
	return ctx.ReconcileFile(
		opt.OutputFileArgs.File,
		func(reconciler reconciling.Reconciler) (*reconciling.Result, error) {
			return reconciler.InsertRecord(date, lines)
		},
	)
}
