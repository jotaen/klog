package cli

import (
	"github.com/jotaen/klog/klog/app"
	"github.com/jotaen/klog/klog/app/cli/util"
)

type Goto struct {
	util.OutputFileArgs
}

func (opt *Goto) Run(ctx app.Context) app.Error {
	target, rErr := ctx.RetrieveTargetFile(opt.File)
	if rErr != nil {
		return rErr
	}

	hasSucceeded := false
	for _, c := range ctx.FileExplorers() {
		c.Args = append(c.Args, target.Location())
		cErr := ctx.Execute(c)
		if cErr != nil {
			continue
		}
		hasSucceeded = true
		break
	}

	if !hasSucceeded {
		return app.NewError(
			"Failed to open file browser",
			"Opening a file browser doesn’t seem possible on your system.",
			nil,
		)
	}
	return nil
}
