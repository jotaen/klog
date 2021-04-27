package cli

import (
	. "klog"
	"klog/app"
	"klog/app/cli/lib"
	"klog/parser"
	"klog/parser/parsing"
)

type Create struct {
	Template    string   `name:"template" hidden help:"The name of the template to instantiate"`
	ShouldTotal Duration `name:"should" help:"The should-total of the record"`
	lib.AtDateArgs
	lib.NoStyleArgs
	lib.OutputFileArgs
}

func (opt *Create) Run(ctx app.Context) error {
	opt.NoStyleArgs.SetGlobalState()
	date := opt.AtDate(ctx.Now())
	lines, err := func() ([]parsing.Text, error) {
		if opt.Template != "" {
			return ctx.InstantiateTemplate(opt.Template)
		}
		headline := opt.AtDate(ctx.Now()).ToString()
		if opt.ShouldTotal != nil {
			headline += " (" + opt.ShouldTotal.ToString() + "!)"
		}
		return []parsing.Text{
			{headline, 0},
		}, nil
	}()
	if err != nil {
		return err
	}
	return applyReconciler(
		opt.OutputFileArgs,
		ctx,
		func(pr *parser.ParseResult) (*parser.ReconcileResult, error) {
			reconciler := parser.NewBlockReconciler(pr, date)
			return reconciler.InsertBlock(lines)
		},
	)
}
