package cli

import (
	"errors"
	"fmt"
	"github.com/alecthomas/kong"
	"klog"
	"klog/app"
	. "klog/lib/jotaen/tf"
	"klog/parser/engine"
	"klog/service"
	"reflect"
	"strings"
)

type cli struct {
	Print  Print  `cmd group:"Evaluate" help:"Pretty-print records"`
	Total  Total  `cmd group:"Evaluate" help:"Evaluate the total time"`
	Eval   Eval   `cmd group:"Evaluate" hidden`
	Report Report `cmd group:"Evaluate" help:"Print a calendar report summarising all days"`
	Tags   Tags   `cmd group:"Evaluate" help:"Print total times aggregated by tags"`

	Append Append `cmd group:"Manipulate" hidden help:"Appends a new record to a file (based on templates)"`

	Bookmark Bookmark `cmd group:"Misc" help:"Set a default file that klog reads from"`
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
		kong.UsageOnError(),
		func() kong.Option {
			datePrototype, _ := klog.NewDate(1, 1, 1)
			return kong.TypeMapper(reflect.TypeOf(&datePrototype).Elem(), dateDecoder())
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

func prettifyError(err error) error {
	switch e := err.(type) {
	case engine.Errors:
		message := ""
		INDENT := "    "
		for _, e := range e.Get() {
			message += fmt.Sprintf(
				Style{Background: "160", Color: "015"}.Format(" ERROR in line %d: "),
				e.Context().LineNumber,
			) + "\n"
			message += fmt.Sprintf(
				Style{Color: "247"}.Format(INDENT+"%s"),
				string(e.Context().Value),
			) + "\n"
			message += fmt.Sprintf(
				Style{Color: "160"}.Format(INDENT+"%s%s"),
				strings.Repeat(" ", e.Position()), strings.Repeat("^", e.Length()),
			) + "\n"
			message += fmt.Sprintf(
				Style{Color: "227"}.Format(INDENT+"%s"),
				strings.Join(breakLines(e.Message(), 60), "\n"+INDENT),
			) + "\n\n"
		}
		return errors.New(message)
	case app.Error:
		return errors.New("Error: " + e.Error() + "\n" + e.Help())
	}
	return errors.New("Error: " + err.Error())
}

func prettifyWarnings(ws []service.Warning) string {
	result := ""
	for _, w := range ws {
		result += Style{Background: "227", Color: "000"}.Format(" WARNING ")
		result += " "
		result += Style{Color: "227"}.Format(w.Date.ToString() + ": " + w.Message)
		result += "\n"
	}
	return result
}

func breakLines(text string, maxLength int) []string {
	SPACE := " "
	words := strings.Split(text, SPACE)
	lines := []string{""}
	for i, w := range words {
		lastLine := lines[len(lines)-1]
		isLastWord := i == len(words)-1
		if !isLastWord && len(lastLine)+len(words[i+1]) > maxLength {
			lines = append(lines, "")
		}
		lines[len(lines)-1] += w + SPACE
	}
	return lines
}
