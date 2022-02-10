/*
Package klog is the entry point of the command line tool.
*/
package klog

import (
	"errors"
	"github.com/alecthomas/kong"
	"github.com/jotaen/klog/src"
	"github.com/jotaen/klog/src/app"
	"github.com/jotaen/klog/src/app/cli"
	"github.com/jotaen/klog/src/app/cli/lib"
	"github.com/jotaen/klog/src/service/period"
	"reflect"
	"strings"
)

func Run(homeDir string, meta app.Meta, isDebug bool, args []string) (int, error) {
	ctx := app.NewContext(homeDir, meta, lib.NewCliSerialiser())
	kongApp, nErr := kong.New(
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
			shouldTotalPrototype := klog.NewShouldTotal(0, 0)
			return kong.TypeMapper(reflect.TypeOf(&shouldTotalPrototype).Elem(), shouldTotalDecoder())
		}(),
		func() kong.Option {
			someSinceDate, _ := klog.NewDate(1, 1, 1)
			someUntilDate, _ := klog.NewDate(2, 2, 2)
			p := period.NewPeriod(someSinceDate, someUntilDate)
			return kong.TypeMapper(reflect.TypeOf(&p).Elem(), periodDecoder())
		}(),
		kong.ConfigureHelp(kong.HelpOptions{
			Compact: true,
		}),
	)
	if nErr != nil {
		return -1, nErr
	}

	kongCtx, cErr := kongApp.Parse(args)
	if cErr != nil {
		return -1, cErr
	}
	kongCtx.BindTo(ctx, (*app.Context)(nil))

	rErr := kongCtx.Run()
	if rErr != nil {
		ctx.Print(lib.PrettifyError(rErr, isDebug).Error() + "\n")
		if appErr, isAppError := rErr.(app.Error); isAppError {
			return int(appErr.Code()), nil
		} else {
			return -1, rErr
		}
	}
	return 0, nil
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

func shouldTotalDecoder() kong.MapperFunc {
	return func(ctx *kong.DecodeContext, target reflect.Value) error {
		var value string
		if err := ctx.Scan.PopValueInto("should", &value); err != nil {
			return err
		}
		if value == "" {
			return errors.New("Please provide a valid should-total duration")
		}
		valueAsDuration := strings.TrimSuffix(value, "!")
		duration, err := klog.NewDurationFromString(valueAsDuration)
		if err != nil {
			return errors.New("`" + value + "` is not a valid should total")
		}
		shouldTotal := klog.NewShouldTotal(0, duration.InMinutes())
		target.Set(reflect.ValueOf(shouldTotal))
		return nil
	}
}

func periodDecoder() kong.MapperFunc {
	return func(ctx *kong.DecodeContext, target reflect.Value) error {
		var value string
		if err := ctx.Scan.PopValueInto("period", &value); err != nil {
			return err
		}
		if value == "" {
			return errors.New("Please provide a valid period")
		}
		p := func() period.Period {
			yearPeriod, yErr := period.NewYearFromString(value)
			if yErr == nil {
				return yearPeriod.Period()
			}
			monthPeriod, mErr := period.NewMonthFromString(value)
			if mErr == nil {
				return monthPeriod.Period()
			}
			return nil
		}()
		if p == nil {
			return errors.New("`" + value + "` is not a valid period")
		}
		target.Set(reflect.ValueOf(p))
		return nil
	}
}
