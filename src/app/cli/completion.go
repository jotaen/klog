package cli

import (
	"github.com/jotaen/klog/src/app"
)

type Completion struct{}

func (c *Completion) Help() string {
	return "Paste the returned code into your shell initialization file, e.g. `~/.bashrc` for Bash."
}

func (c *Completion) Run(ctx app.Context) error {
	completion, err := ctx.Completion()
	if err != nil {
		return err
	}
	ctx.Print(completion)
	return nil
}
