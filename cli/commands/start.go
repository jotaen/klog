package commands

import (
	"klog/cli"
	"klog/datetime"
	"time"
)

var Start cli.Command

func init() {
	Start = cli.Command{
		Name:        "start",
		Alias:       []string{},
		Description: "Create a new entry",
		Main:        start,
	}
}

func start(env cli.Environment, args []string) int {
	now := time.Now()
	today, _ := datetime.CreateDateFromTime(now)
	wd, err := env.Store.Get(today)
	if err != nil {
		// todo create new
		return 182763
	}
	nowTime, _ := datetime.CreateTimeFromTime(now)
	openTimeRange, _ := datetime.CreateTimeRange(nowTime, nil)
	wd.AddRange(openTimeRange)
	env.Store.Save(wd)
	return cli.OK
}
