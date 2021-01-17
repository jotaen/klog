package cli

import (
	"fmt"
	"github.com/alecthomas/kong"
	"klog/app"
	"klog/record"
	"reflect"
)

var cli struct {
	Print  Print  `cmd help:"Show records in a file"`
	Total  Total  `cmd help:"Sum up the total time"`
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
		func() kong.Option {
			datePrototype, _ := record.NewDate(1, 1, 1)
			return kong.TypeMapper(reflect.TypeOf(&datePrototype).Elem(), dateDecoder())
		}(),
	)
	err = args.Run(ctx)
	if err != nil {
		fmt.Println(err)
		return -1
	}
	return 0
}

func dateDecoder() kong.MapperFunc {
	return func(ctx *kong.DecodeContext, target reflect.Value) error {
		var value string
		if err := ctx.Scan.PopValueInto("date", &value); err != nil {
			return err
		}
		d, err := record.NewDateFromString(value)
		if err != nil {
			return err
		}
		target.Set(reflect.ValueOf(d))
		return nil
	}
}
