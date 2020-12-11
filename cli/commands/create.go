package commands

import (
	"github.com/akamensky/argparse"
	"klog/cli"
	"klog/datetime"
	"klog/workday"
	"time"
)

func Create(env cli.Environment, args []string) int {
	opts, err := parseArgs(args)
	if err != nil {
		return cli.INVALID_CLI_ARGS
	}
	wd := workday.Create(opts.date)
	env.Store.Save(wd)
	return cli.OK
}

type opts struct {
	date datetime.Date
}

func parseArgs(args []string) (opts, error) {
	argParser := argparse.NewParser("create", "")
	dateArg := argParser.String("d", "date", &argparse.Options{
		Required: false,
		Default:  "today",
	})
	err := argParser.Parse(args)
	opts := opts{}
	if err != nil {
		return opts, err
	}
	if *dateArg == "" || *dateArg == "today" {
		date, _ := datetime.CreateDateFromTime(time.Now())
		opts.date = date
		return opts, nil
	} else {
		date, err := datetime.CreateDateFromString(*dateArg)
		if err != nil {
			return opts, err
		}
		opts.date = date
		return opts, nil
	}
}
