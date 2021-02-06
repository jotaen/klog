package cli

import (
	"klog/app"
	systray "klog/app/mac_widget"
)

type Widget struct {
	Detach bool `name:"detach" help:"Detach the widget from the command line"`
}

func (args *Widget) Run(ctx app.Context) error {
	if !args.Detach {
		ctx.Print("If you would like to run the widget on its own\n")
		ctx.Print("run this command again with --detach\n")
	}
	systray.Run(args.Detach)
	return nil
}
