/*
Package klog is the entry point of the command line tool.
*/
package klog

import (
	"github.com/alecthomas/kong"
	"github.com/jotaen/klog/klog"
	"github.com/jotaen/klog/klog/app"
	"github.com/jotaen/klog/klog/app/cli"
	"github.com/jotaen/klog/klog/app/cli/lib"
	"github.com/jotaen/klog/klog/service"
	"github.com/jotaen/klog/klog/service/period"
	kongcompletion "github.com/jotaen/kong-completion"
	"reflect"
)

func Run(homeDir app.File, meta app.Meta, config app.Config, args []string) error {
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
		func() kong.Option {
			f, _ := service.NewRounding(30)
			return kong.TypeMapper(reflect.TypeOf(&f).Elem(), roundingDecoder())
		}(),
		func() kong.Option {
			t := klog.NewTagOrPanic("test", "")
			return kong.TypeMapper(reflect.TypeOf(&t).Elem(), tagDecoder())
		}(),
		func() kong.Option {
			s, _ := klog.NewRecordSummary("test")
			return kong.TypeMapper(reflect.TypeOf(&s).Elem(), recordSummaryDecoder())
		}(),
		func() kong.Option {
			s, _ := klog.NewEntrySummary("test")
			return kong.TypeMapper(reflect.TypeOf(&s).Elem(), entrySummaryDecoder())
		}(),
		kong.ConfigureHelp(kong.HelpOptions{
			Compact:             true,
			NoExpandSubcommands: true,
		}),
	)
	if nErr != nil {
		return nErr
	}

	ctx := app.NewContext(
		homeDir,
		meta,
		lib.CliSerialiser{},
		config,
	)

	// When klog is invoked by shell completion (specifically, when the
	// bash-specific COMP_LINE environment variable is set), the
	// kongplete.Complete function generates a list of possible completions,
	// prints them one per line to stdout, and then exits the program early.
	kongcompletion.Register(
		kongApp,
		kongcompletion.WithPredictors(CompletionPredictors(ctx)),
		kongcompletion.WithFlagOverrides(lib.FilterArgsCompletionOverrides),
	)

	kongCtx, cErr := kongApp.Parse(args)
	if cErr != nil {
		return cErr
	}
	kongCtx.BindTo(ctx, (*app.Context)(nil))

	return kongCtx.Run()
}
