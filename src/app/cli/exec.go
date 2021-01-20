package cli

import (
	"errors"
	"fmt"
	"github.com/alecthomas/kong"
	"klog"
	"klog/app"
	"reflect"
	"time"
)

var cli struct {
	Print  Print  `cmd help:"Show records in a file"`
	Total  Total  `cmd help:"Sum up the total time"`
	Diff   Diff   `cmd help:"Show diff between total and should time"`
	Follow Follow `cmd help:"Follow today’s record live"`
	Widget Widget `cmd help:"Launch menu bar widget (MacOS only)"`
}

func Execute() int {
	ctx, err := app.NewContextFromEnv()
	if err != nil {
		fmt.Println("Failed to initialise application. Error:")
		fmt.Println(err)
		return -1
	}
	args := kong.Parse(
		&cli,
		kong.Name("klog"),
		kong.Description("klog time tracking: command line app for interacting with `.klg` files."),
		kong.UsageOnError(),
		func() kong.Option {
			datePrototype, _ := src.NewDate(1, 1, 1)
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
		if value == "today" || value == "yesterday" {
			now := time.Now()
			if value == "yesterday" {
				now = time.Now().AddDate(0, 0, -1)
			}
			value = fmt.Sprintf("%04d-%02d-%02d", now.Year(), now.Month(), now.Day())
		}

		d, err := src.NewDateFromString(value)
		if err != nil {
			return errors.New("`" + value + "` is not a valid date")
		}
		target.Set(reflect.ValueOf(d))
		return nil
	}
}
