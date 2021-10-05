package cli

import (
	"github.com/jotaen/klog/src/app"
	"github.com/jotaen/klog/src/app/cli/lib"
)

type Edit struct {
	lib.OutputFileArgs
	lib.QuietArgs
}

func (opt *Edit) Run(ctx app.Context) error {
	printHint := (func() func(string) {
		if opt.Quiet {
			return func(string) {}
		}
		return ctx.Print
	})()
	return ctx.OpenInEditor(opt.File, printHint)
}
