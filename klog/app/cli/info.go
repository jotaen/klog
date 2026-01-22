package cli

import (
	"github.com/jotaen/klog/klog/app"
)

type Info struct {
	Spec      bool `name:"spec" help:"Print the .klg file format specification."`
	License   bool `name:"license" help:"Print license / copyright information."`
	About     bool `name:"about" help:"Print meta information about klog."`
	Filtering bool `name:"filtering" help:"Print documentation for using filter expressions."`
}

func (opt *Info) Run(ctx app.Context) app.Error {
	if opt.Spec {
		ctx.Print(ctx.Meta().Specification + "\n")
	} else if opt.License {
		ctx.Print(ctx.Meta().License + "\n")
	} else if opt.About {
		ctx.Print(INTRO_SUMMARY)
	} else if opt.Filtering {
		ctx.Print(`klog filter expressions are a generic, predicate-based language for filtering data for evaluation purposes.

`) // TODO write documentation
	} else {
		return app.NewErrorWithCode(
			app.GENERAL_ERROR,
			"No flag specified",
			"Run with `--help` for more info",
			nil,
		)
	}
	return nil
}
