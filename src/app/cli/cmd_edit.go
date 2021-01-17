package cli

import (
	"klog/app"
)

var Edit Command

func init() {
	Edit = Command{
		Name:        "edit",
		Description: "Open file in editor",
		Main:        edit,
	}
}

func edit(ctx app.Context, args []string) int {
	err := ctx.OpenInEditor()
	if err != nil {
		return EXECUTION_FAILED
	}
	return OK
}
