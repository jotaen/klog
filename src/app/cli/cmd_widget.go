package cli

import (
	"klog/app"
	systray "klog/app/mac_widget"
)

type Widget struct {
	Detach bool `name:"detach" help:"Detach the widget from the command line"`
}

func (opt *Widget) Run(ctx app.Context) error {
	if !opt.Detach {
		ctx.Print("If you would like to run the widget on its own\n")
		ctx.Print("run this command again with --detach\n")
	}
	systray.Run(opt.Detach)
	return nil
}
