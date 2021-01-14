package commands

import (
	"github.com/akamensky/argparse"
	"klog/app"
	"klog/app/cli"
	"klog/datetime"
	"klog/project"
	"klog/record"
	"time"
)

var Create cli.Command

func init() {
	Create = cli.Command{
		Name:        "create",
		Alias:       []string{"new"},
		Description: "Create a new entry",
		Main:        create,
	}
}

func create(env app.Environment, project project.Project, args []string) int {
	opts, err := parseArgs(args)
	if err != nil {
		return cli.INVALID_CLI_ARGS
	}
	wd := record.NewRecord(opts.date)
	project.Save(wd)
	return cli.OK
}

type opts struct {
	date datetime.Date
}

func parseArgs(args []string) (opts, error) {
	argParser := argparse.NewParser(Create.Name, Create.Description)
	dateArg := argParser.String("d", "date", &argparse.Options{
		Required: false,
		Default:  "today",
		Help:     "Provide a date (format: YYYY-MM-DD or `today`)",
	})
	err := argParser.Parse(args)
	opts := opts{}
	if err != nil {
		return opts, err
	}
	if *dateArg == "" || *dateArg == "today" {
		date, _ := datetime.NewDateFromTime(time.Now())
		opts.date = date
		return opts, nil
	} else {
		date, err := datetime.NewDateFromString(*dateArg)
		if err != nil {
			return opts, err
		}
		opts.date = date
		return opts, nil
	}
}
