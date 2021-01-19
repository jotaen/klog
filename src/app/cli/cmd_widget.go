package cli

import (
	"klog/app"
	systray "klog/app/mac_widget"
)

type Widget struct {
	Start bool   `name:"start" help:"Launch widget"`
	File  string `short:"f" name:"file" help:"Specify file"`
}

func (args *Widget) Run(ctx *app.Context) error {
	if args.Start {
		systray.Launch()
	}
	return nil
}
