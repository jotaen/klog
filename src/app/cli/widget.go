package cli

import (
	"github.com/jotaen/klog/src/app"
	systray "github.com/jotaen/klog/src/app/mac_widget"
)

type Widget struct {
	Detach bool `name:"detach" help:"Detach the widget from the command line"`
}

func (opt *Widget) Help() string {
	return `MacOS only! The widget is a small application that shows up in the menubar (systray),
next to the system clock.

The widget reads from the bookmarked file and presents a brief summary of the data.

Note that this is an experimental feature, it might be discontinued.`
}

func (opt *Widget) Run(ctx app.Context) error {
	if !systray.IsWidgetAvailable() {
		return app.NewError(
			"Cannot start widget",
			"The widget is currently only supported on MacOS.",
			nil,
		)
	}
	if !opt.Detach {
		ctx.Print("If you would like to run the widget on its own\n")
		ctx.Print("run this command again with --detach\n")
	}
	systray.Run(opt.Detach)
	return nil
}
