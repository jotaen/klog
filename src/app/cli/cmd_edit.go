package cli

import (
	"errors"
	"klog/app"
)

type Edit struct {
	FileArgs
}

func (args *Edit) Run(ctx *app.Context) error {
	err := ctx.OpenInEditor()
	if err != nil {
		return errors.New("EXECUTION_FAILED")
	}
	return nil
}
