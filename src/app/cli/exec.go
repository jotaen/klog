package cli

import (
	"errors"
	"fmt"
	"github.com/alecthomas/kong"
	"klog"
	"klog/app"
	"reflect"
	"regexp"
	"strconv"
	"strings"
)

type cli struct {
	Print  Print  `cmd group:"Evaluate" help:"Pretty-print records"`
	Total  Total  `cmd group:"Evaluate" help:"Evaluate the total time"`
	Report Report `cmd group:"Evaluate" help:"Print a calendar report summarising all days"`
	Tags   Tags   `cmd group:"Evaluate" help:"Print total times aggregated by tags"`
	Now    Now    `cmd group:"Evaluate" help:"Evaluate today’s record (including potential open ranges)"`

	Append Append `cmd group:"Manipulate" hidden help:"Appends a new record to a file (based on templates)"`

	Bookmark Bookmark `cmd group:"Misc" help:"Default file that klog reads from"`
	Json     Json     `cmd group:"Misc" help:"Convert records to JSON"`
	Widget   Widget   `cmd group:"Misc" help:"Start menu bar widget (MacOS only)"`
	Version  Version  `cmd group:"Misc" help:"Print version info and check for updates"`
}

func Execute() int {
	ctx, err := app.NewContextFromEnv()
	if err != nil {
		fmt.Println("Failed to initialise application. Error:")
		fmt.Println(err)
		return -1
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
			period := Period{}
			return kong.TypeMapper(reflect.TypeOf(&period).Elem(), periodDecoder())
		}(),
		kong.ConfigureHelp(kong.HelpOptions{
			Compact: true,
		}),
	)
	cliApp.BindTo(ctx, (*app.Context)(nil))
	err = cliApp.Run(&ctx)
	if err != nil {
		fmt.Println(prettifyError(err))
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

var periodPattern = regexp.MustCompile(`^\d{4}(-\d{2})?$`)

func periodDecoder() kong.MapperFunc {
	return func(ctx *kong.DecodeContext, target reflect.Value) error {
		var value string
		if err := ctx.Scan.PopValueInto("period", &value); err != nil {
			return err
		}
		if value == "" || !periodPattern.MatchString(value) {
			return errors.New("Please provide a valid period")
		}
		parts := strings.Split(value, "-")
		year, _ := strconv.Atoi(parts[0])
		monthStart := 1
		monthEnd := 12
		if len(parts) == 2 {
			monthStart, _ = strconv.Atoi(parts[1])
			monthEnd = monthStart
		}
		start, _ := klog.NewDate(year, monthStart, 1)
		end, _ := klog.NewDate(year, monthEnd, 28)
		for true {
			next := end.PlusDays(1)
			if next.Month() != end.Month() {
				break
			}
			end = next
		}
		target.Set(reflect.ValueOf(Period{
			since: start,
			until: end,
		}))
		return nil
	}
}

type Period struct {
	since klog.Date
	until klog.Date
}
