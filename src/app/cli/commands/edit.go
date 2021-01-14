package commands

import (
	"klog/app"
	"klog/app/cli"
)

var Edit cli.Command

func init() {
	Edit = cli.Command{
		Name:        "edit",
		Description: "Open file in editor",
		Main:        edit,
	}
}

func edit(service app.Service, args []string) int {
	err := service.OpenInEditor()
	if err != nil {
		return cli.EXECUTION_FAILED
	}
	return cli.OK
}
