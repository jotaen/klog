package cli

import (
	"fmt"
	"klog/app"
	systray "klog/app/mac_widget"
)

type Widget struct {
	File   string `short:"f" name:"file" help:"Which file to show in the widget"`
	Detach bool   `name:"detach" help:"Detach the widget from the cli"`
}

func (args *Widget) Run(ctx *app.Context) error {
	if args.File != "" {
		err := ctx.SetBookmark(args.File)
		if err != nil {
			return err
		}
	}
	if !args.Detach {
		fmt.Println("If you would like to run the widget on its own")
		fmt.Println("run again with --detach")
	}
	systray.Run(args.Detach)
	return nil
}
