/*
Package klog is the entry point of the command line tool.
*/
package klog

import (
	"errors"
	"reflect"

	"github.com/alecthomas/kong"
	"github.com/jotaen/klog/klog"
	"github.com/jotaen/klog/klog/app"
	"github.com/jotaen/klog/klog/app/cli"
	"github.com/jotaen/klog/klog/app/cli/prettify"
	"github.com/jotaen/klog/klog/service"
	"github.com/jotaen/klog/klog/service/period"
	tf "github.com/jotaen/klog/lib/terminalformat"
	kongcompletion "github.com/jotaen/kong-completion"
)

func Run(homeDir app.File, meta app.Meta, config app.Config, args []string) (int, error) {
	kongApp, nErr := kong.New(
		&cli.Cli{},
		kong.Name("klog"),
		kong.Description((&cli.Default{}).Help()),
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
		func() kong.Option {
			t := service.ENTRY_TYPE_DURATION
			return kong.TypeMapper(reflect.TypeOf(&t).Elem(), entryTypeDecoder())
		}(),
		kong.ConfigureHelp(kong.HelpOptions{
			Compact:             true,
			NoExpandSubcommands: true,
			WrapUpperBound:      prettify.LINE_LENGTH,
		}),
	)
	if nErr != nil {
		// This code branch is not expected to be invoked in practice. If it were to
		// happen, that most likely indicates a bug in the app setup.
		return app.GENERAL_ERROR.ToInt(), errors.New("Internal error: " + nErr.Error())
	}

	styler := tf.NewStyler(config.ColourScheme.Value())
	ctx := app.NewContext(
		homeDir,
		meta,
		styler,
		config,
	)

	// When klog is invoked by shell completion (specifically, when the
	// bash-specific COMP_LINE environment variable is set), the
	// kongplete.Complete function generates a list of possible completions,
	// prints them one per line to stdout, and then exits the program early.
	kongcompletion.Register(
		kongApp,
		kongcompletion.WithPredictors(CompletionPredictors(ctx)),
	)

	kongCtx, cErr := kongApp.Parse(args)
	if cErr != nil {
		return app.GENERAL_ERROR.ToInt(), errors.New("Invocation error: " + cErr.Error())
	}
	kongCtx.BindTo(ctx, (*app.Context)(nil))

	rErr := kongCtx.Run()
	if rErr == nil {
		return 0, nil
	}

	appError := app.NewError("", "", nil)
	if errors.As(rErr, &appError) {
		parserErrors := app.NewParserErrors(nil)
		switch {
		case errors.As(appError, &parserErrors):
			return parserErrors.Code().ToInt(), prettify.PrettifyParsingError(parserErrors, styler)
		}
		return appError.Code().ToInt(), prettify.PrettifyAppError(appError, config.IsDebug.Value())
	}

	// This is just a fallback, which is not expected to occur in practice.
	return app.GENERAL_ERROR.ToInt(), errors.New("Error: " + rErr.Error())
}
