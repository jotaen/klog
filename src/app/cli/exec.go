package cli

import (
	"errors"
	"fmt"
	"github.com/alecthomas/kong"
	"klog"
	"klog/app"
	"reflect"
)

type cli struct {
	Print   Print   `cmd help:"Pretty-print records"`
	Eval    Eval    `cmd help:"Evaluate records"`
	Widget  Widget  `cmd help:"Start menu bar widget (MacOS only)"`
	Version Version `cmd help:"Print version info and check for updates"`
}

func Execute() int {
	ctx, err := app.NewContextFromEnv()
	if err != nil {
		fmt.Println("Failed to initialise application. Error:")
		fmt.Println(err)
		return -1
	}
	args := kong.Parse(
		&cli{},
		kong.Name("klog"),
		kong.Description("klog time tracking: command line app for interacting with `.klg` files."),
		kong.UsageOnError(),
		func() kong.Option {
			datePrototype, _ := klog.NewDate(1, 1, 1)
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
		if value == "" {
			return errors.New("please provide a valid date")
		}
		d, err := klog.NewDateFromString(value)
		if err != nil {
			return errors.New("`" + value + "` is not a valid date")
		}
		target.Set(reflect.ValueOf(d))
		return nil
	}
}
