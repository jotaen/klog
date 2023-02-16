package cli

import (
	"github.com/jotaen/klog/klog/app"
	"github.com/jotaen/klog/klog/app/cli/lib"
	"github.com/jotaen/klog/klog/app/cli/lib/command"
)

type Edit struct {
	lib.OutputFileArgs
	lib.QuietArgs
}

var hint = "You can specify your preferred editor via the $EDITOR environment variable, or the klog config file."

func (opt *Edit) Run(ctx app.Context) app.Error {
	target, err := ctx.RetrieveTargetFile(opt.File)
	if err != nil {
		return err
	}

	explicitEditor, autoEditors := ctx.Editors()

	if explicitEditor != "" {
		c, cErr := command.NewFromString(explicitEditor)
		if cErr != nil {
			return app.NewError(
				"Invalid editor setting",
				"Please check the value for invalid syntax: "+explicitEditor,
				cErr,
			)
		}
		c.Args = append(c.Args, target.Path())
		rErr := ctx.Execute(c)
		if rErr != nil {
			return app.NewError(
				"Cannot open preferred editor",
				"Editor command was: "+explicitEditor+"\nNote that if your editor path contains spaces, you have to quote it.",
				nil,
			)
		}
	} else {
		hasSucceeded := false
		for _, c := range autoEditors {
			c.Args = append(c.Args, target.Path())
			rErr := ctx.Execute(c)
			if rErr == nil {
				hasSucceeded = true
				break
			}
		}

		if !hasSucceeded {
			return app.NewError(
				"Cannot open any editor",
				hint,
				nil,
			)
		}

		if !opt.Quiet {
			ctx.Print(hint + "\n")
		}
	}

	return nil
}
