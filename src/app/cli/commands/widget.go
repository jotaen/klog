package commands

import (
	"klog/app"
	"klog/app/cli"
	systray "klog/app/mac_widget"
)

var Widget cli.Command

func init() {
	Widget = cli.Command{
		Name:        "widget",
		Description: "Launch widget in systray",
		Main:        widget,
	}
}

func widget(_ app.Service, _ []string) int {
	systray.Launch()
	return cli.OK
}
