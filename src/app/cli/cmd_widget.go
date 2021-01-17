package cli

import (
	"klog/app"
	systray "klog/app/mac_widget"
)

var Widget Command

func init() {
	Widget = Command{
		Name:        "widget",
		Description: "Launch widget in systray",
		Main:        widget,
	}
}

func widget(_ app.Context, _ []string) int {
	systray.Launch()
	return OK
}
