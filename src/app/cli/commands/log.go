package commands

import (
	"klog/app"
	"klog/app/cli"
	"klog/datetime"
	"klog/project"
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

func log(env app.Environment, project project.Project, args []string) int {
	now := time.Now()
	today, _ := datetime.NewDateFromTime(now)
	wd, err := project.Get(today)
	if err != nil {
		// todo create new
		return 182763
	}
	value := strings.Join(args[:], "")
	duration, _ := datetime.NewDurationFromString(value)
	wd.AddDuration(duration)
	project.Save(wd)
	return cli.OK
}
