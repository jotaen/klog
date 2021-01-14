package commands

import (
	"github.com/akamensky/argparse"
	"klog/app"
	"klog/app/cli"
	"klog/datetime"
	"klog/record"
	"time"
)

var Create cli.Command

func init() {
	Create = cli.Command{
		Name:        "create",
		Description: "Create a new entry",
		Main:        create,
	}
}

func create(service app.Service, args []string) int {
	opts, err := parseArgs(args)
	if err != nil {
		return cli.INVALID_CLI_ARGS
	}
	rs := service.Input()
	r := record.NewRecord(opts.date)
	newRecords := append(rs, r)
	service.Save(newRecords)
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
