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
	ShouldTotal Duration `name:"should" help:"A should total property"`
	lib.AtDateArgs
	lib.OutputFileArgs
}

func (opt *Create) Run(ctx app.Context) error {
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
	return reconcile(
		opt.OutputFileArgs,
		ctx,
		func(pr *parser.ParseResult) (*parser.Reconciler, error) {
			return parser.NewBlockReconciler(pr, func(r1 Record, r2 Record) bool {
				return date.IsAfterOrEqual(r1.Date()) && r2.Date().IsAfterOrEqual(date)
			})
		},
		func(r *parser.Reconciler) (Record, string, error) {
			return r.AddBlock(lines)
		},
	)
}
