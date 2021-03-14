package cli

import (
	"errors"
	. "klog"
	"klog/app"
	"klog/app/cli/lib"
	"klog/parser"
	"klog/parser/parsing"
)

type Create struct {
	Template string `required name:"template" short:"t" help:"The name of the template to instantiate"`
	lib.AtDateArgs
	lib.OutputFileArgs
}

func (opt *Create) Run(ctx app.Context) error {
	date := opt.AtDate(ctx.Now())
	lines, err := func() ([]parsing.Text, error) {
		if opt.Template != "" {
			return ctx.InstantiateTemplate(opt.Template)
		}
		return []parsing.Text{
			{opt.AtDate(ctx.Now()).ToString(), 0},
		}, nil
	}()
	if err != nil {
		return err
	}
	return reconcile(
		opt.OutputFileArgs,
		ctx,
		errors.New("No eligible record at date "+date.ToString()),
		func(r Record) bool { return date.IsAfterOrEqual(r.Date()) },
		func(r *parser.Reconciler) (Record, string, error) {
			return r.AddBlock(lines)
		},
	)
}
