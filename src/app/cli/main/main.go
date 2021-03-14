package main

import (
	"fmt"
	"github.com/alecthomas/kong"
	"klog"
	"klog/app"
	"klog/app/cli/lib"
	"os"
	"reflect"
)

func main() {
	ctx, err := app.NewContextFromEnv()
	if err != nil {
		fmt.Println("Failed to initialise application. Error:")
		fmt.Println(err)
		os.Exit(-1)
	}
	cliApp := kong.Parse(
		&cli{},
		kong.Name("klog"),
		kong.Description(
			"klog time tracking: command line app for interacting with `.klg` files."+
				"\n\nRead the documentation at https://klog.jotaen.net",
		),
		func() kong.Option {
			datePrototype, _ := klog.NewDate(1, 1, 1)
			return kong.TypeMapper(reflect.TypeOf(&datePrototype).Elem(), dateDecoder())
		}(),
		func() kong.Option {
			timePrototype, _ := klog.NewTime(0, 0)
			return kong.TypeMapper(reflect.TypeOf(&timePrototype).Elem(), timeDecoder())
		}(),
		func() kong.Option {
			durationPrototype := klog.NewDuration(0, 0)
			return kong.TypeMapper(reflect.TypeOf(&durationPrototype).Elem(), durationDecoder())
		}(),
		func() kong.Option {
			period := lib.Period{}
			return kong.TypeMapper(reflect.TypeOf(&period).Elem(), periodDecoder())
		}(),
		kong.ConfigureHelp(kong.HelpOptions{
			Compact: true,
		}),
	)
	cliApp.BindTo(ctx, (*app.Context)(nil))
	err = cliApp.Run(&ctx)
	if err != nil {
		fmt.Println(lib.PrettifyError(err))
		os.Exit(-1)
	}
	os.Exit(0)
}
