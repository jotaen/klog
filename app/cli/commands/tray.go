package commands

import (
	"klog/app/cli"
	systray "klog/app/tray"
)

var Tray cli.Command

func init() {
	Tray = cli.Command{
		Name:        "tray",
		Alias:       []string{},
		Description: "Launch widget in systray",
		Main:        tray,
	}
}

func tray(env cli.Environment, args []string) int {
	systray.Start()
	return cli.OK
}
