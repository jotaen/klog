package cli

import (
	"klog/app"
	systray "klog/app/mac_widget"
)

type Widget struct {
	File   string `short:"f" name:"file" help:"Set the file to read from"`
	Detach bool   `name:"detach" help:"Detach the widget from the command line"`
}

func (args *Widget) Run(ctx app.Context) error {
	if args.File != "" {
		ctx.Print("Set file " + args.File)
		err := ctx.SetBookmark(args.File)
		return err
	}
	if !args.Detach {
		ctx.Print("If you would like to run the widget on its own")
		ctx.Print("run again with --detach")
	}
	systray.Run(args.Detach)
	return nil
}
