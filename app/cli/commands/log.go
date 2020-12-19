package commands

import (
	"klog/app/cli"
	"klog/datetime"
	"strings"
	"time"
)

var Log cli.Command

func init() {
	Log = cli.Command{
		Name:        "log",
		Alias:       []string{},
		Description: "Create a new entry",
		Main:        log,
	}
}

func log(env cli.Environment, args []string) int {
	now := time.Now()
	today, _ := datetime.NewDateFromTime(now)
	wd, err := env.Store.Get(today)
	if err != nil {
		// todo create new
		return 182763
	}
	value := strings.Join(args[:], "")
	duration, _ := datetime.NewDurationFromString(value)
	wd.AddDuration(duration)
	env.Store.Save(wd)
	return cli.OK
}
