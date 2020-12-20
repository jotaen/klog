package commands

import (
	"klog/app"
	"klog/app/cli"
	systray "klog/app/tray"
	"klog/project"
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

func tray(env app.Environment, project project.Project, args []string) int {
	systray.Start(env)
	return cli.OK
}
