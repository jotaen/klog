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
