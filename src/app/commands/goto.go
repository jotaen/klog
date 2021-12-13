package commands

import (
	"github.com/jotaen/klog/src/app"
	"github.com/jotaen/klog/src/app/lib"
)

type Goto struct {
	lib.OutputFileArgs
}

func (opt *Goto) Run(ctx app.Context) error {
	return ctx.OpenInFileBrowser(opt.File)
}
