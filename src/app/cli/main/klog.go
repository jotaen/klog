package main

import (
	"errors"
	"fmt"
	"github.com/alecthomas/kong"
	"github.com/jotaen/klog/src"
	"github.com/jotaen/klog/src/app"
	"github.com/jotaen/klog/src/app/cli"
	"github.com/jotaen/klog/src/app/cli/lib"
	"os"
	"reflect"
	"strings"
)

func main() {
	ctx, err := app.NewContextFromEnv(lib.NewCliSerialiser())
	if err != nil {
		fmt.Println("Failed to initialise application. Error:")
		fmt.Println(err)
		os.Exit(-1)
	}
	cliApp := kong.Parse(
		&cli.Cli{},
		kong.Name("klog"),
		kong.Description(cli.DESCRIPTION),
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
		isDebug := false
		if os.Getenv("KLOG_DEBUG") != "" {
			isDebug = true
		}
		fmt.Println(lib.PrettifyError(err, isDebug))
		exitCode := app.GENERAL_ERROR
		if appErr, isAppError := err.(app.Error); isAppError {
			exitCode = appErr.Code()
		}
		os.Exit(exitCode.ToInt())
	}
	os.Exit(0)
}

func dateDecoder() kong.MapperFunc {
	return func(ctx *kong.DecodeContext, target reflect.Value) error {
		var value string
		if err := ctx.Scan.PopValueInto("date", &value); err != nil {
			return err
		}
		if value == "" {
			return errors.New("Please provide a valid date")
		}
		d, err := klog.NewDateFromString(value)
		if err != nil {
			return errors.New("`" + value + "` is not a valid date")
		}
		target.Set(reflect.ValueOf(d))
		return nil
	}
}

func timeDecoder() kong.MapperFunc {
	return func(ctx *kong.DecodeContext, target reflect.Value) error {
		var value string
		if err := ctx.Scan.PopValueInto("time", &value); err != nil {
			return err
		}
		if value == "" {
			return errors.New("Please provide a valid time")
		}
		t, err := klog.NewTimeFromString(value)
		if err != nil {
			return errors.New("`" + value + "` is not a valid time")
		}
		target.Set(reflect.ValueOf(t))
		return nil
	}
}

func durationDecoder() kong.MapperFunc {
	return func(ctx *kong.DecodeContext, target reflect.Value) error {
		var value string
		if err := ctx.Scan.PopValueInto("duration", &value); err != nil {
			return err
		}
		if value == "" {
			return errors.New("Please provide a valid duration")
		}
		value = strings.TrimSuffix(value, "!")
		d, err := klog.NewDurationFromString(value)
		if err != nil {
			return errors.New("`" + value + "` is not a valid duration")
		}
		target.Set(reflect.ValueOf(d))
		return nil
	}
}

func periodDecoder() kong.MapperFunc {
	return func(ctx *kong.DecodeContext, target reflect.Value) error {
		var value string
		if err := ctx.Scan.PopValueInto("period", &value); err != nil {
			return err
		}
		p, err := lib.NewPeriodFromString(value)
		if err != nil {
			return err
		}
		target.Set(reflect.ValueOf(p))
		return nil
	}
}

func fileOrBookmarkDecoder() kong.MapperFunc {
	return func(ctx *kong.DecodeContext, target reflect.Value) error {
		var value string
		if err := ctx.Scan.PopValueInto("fileOrBookmarkName", &value); err != nil {
			return err
		}
		fileOrBookmarkName := app.FileOrBookmarkName(value)
		target.Set(reflect.ValueOf(fileOrBookmarkName))
		return nil
	}
}
