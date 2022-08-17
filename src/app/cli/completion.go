package cli

import (
	"github.com/jotaen/klog/src/app"
)

type Completion struct{}

func (c *Completion) Help() string {
	return "The printed shell code is for instructing your shell to use tab completions for klog. " +
		"Place the code into your shell initialization file, e.g. `~/.bashrc`. " +
		"You can either paste it verbatim, or you source it dynamically via `. <(klog completion)`."
}

func (c *Completion) Run(ctx app.Context) error {
	completion, err := ctx.Completion()
	if err != nil {
		return err
	}
	ctx.Print(completion)
	return nil
}
