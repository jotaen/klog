package commands

import (
	. "github.com/jotaen/klog/src"
	"github.com/jotaen/klog/src/app"
	"github.com/jotaen/klog/src/app/lib"
	"github.com/jotaen/klog/src/parser"
	"github.com/jotaen/klog/src/parser/reconciler"
)

type Create struct {
	Template    string   `name:"template" hidden help:"The name of the template to instantiate"`
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
	lines, err := func() ([]reconciler.Text, error) {
		if opt.Template != "" {
			return ctx.InstantiateTemplate(opt.Template)
		}
		headline := opt.AtDate(ctx.Now()).ToString()
		if opt.ShouldTotal != nil {
			headline += " (" + opt.ShouldTotal.ToString() + "!)"
		}
		return []reconciler.Text{
			{headline, 0},
		}, nil
	}()
	if err != nil {
		return err
	}
	return lib.ReconcilerChain{
		File: opt.OutputFileArgs.File,
		Ctx:  ctx,
	}.Apply(
		func(pr *parser.ParseResult) (*reconciler.ReconcileResult, error) {
			reconciler := reconciler.NewBlockReconciler(pr, date)
			return reconciler.InsertBlock(lines)
		},
	)
}
