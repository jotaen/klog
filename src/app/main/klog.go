/*
Package klog is the entry point of the command line tool.
*/
package klog

import (
	"errors"
	"fmt"
	"github.com/alecthomas/kong"
	"github.com/jotaen/klog/src"
	"github.com/jotaen/klog/src/app"
	"github.com/jotaen/klog/src/app/commands"
	"github.com/jotaen/klog/src/app/lib"
	"os"
	"reflect"
	"strings"
)

func Run(meta app.Meta) app.Error {
	ctx, err := app.NewContextFromEnv(meta, lib.NewCliSerialiser())
	if err != nil {
		fmt.Println("Failed to initialise application. Error:")
		fmt.Println(err)
		os.Exit(-1)
	}
	cliApp := kong.Parse(
		&commands.Cli{},
		kong.Name("klog"),
		kong.Description(commands.DESCRIPTION),
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
		ctx.Print(lib.PrettifyError(err, isDebug).Error() + "\n")
		if appErr, isAppError := err.(app.Error); isAppError {
			return appErr
		} else {
			// Remember that kong errors always panic (such as flag parsing),
			// so we can’t handle them here.
			return app.NewErrorWithCode(
				app.GENERAL_ERROR,
				"An unknown error occurred",
				err.Error(),
				err,
			)
		}
	}
	return nil
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
