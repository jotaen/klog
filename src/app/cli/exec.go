package cli

import (
	"fmt"
	"github.com/alecthomas/kong"
	"klog/app"
)

var cli struct {
	Print  Print  `cmd help:"Show records in a file"`
	Edit   Edit   `cmd help:"Open file in editor"`
	Start  Start  `cmd help:"Start to track"`
	Widget Widget `cmd help:"Launch menu bar widget (MacOS only)"`
}

func Execute() int {
	ctx, err := app.NewContextFromEnv()
	if err != nil {
		return -1
	}
	args := kong.Parse(
		&cli,
		kong.Name("klog"),
		kong.Description("klog time tracking\nCommand line interface for interacting with `.klg` files."),
	)
	err = args.Run(ctx)
	if err != nil {
		fmt.Println(err)
		return -1
	}
	return 0
}
