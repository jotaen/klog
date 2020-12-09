package commands

import (
	"klog/cli"
	"klog/datetime"
	"strings"
	"time"
)

func Log(env cli.Environment, args []string) int {
	now := time.Now()
	today, _ := datetime.CreateDateFromTime(now)
	wd, err := env.Store.Get(today)
	if err != nil {
		// todo create new
		return 182763
	}
	value := strings.Join(args[:], "")
	duration, _ := datetime.CreateDurationFromString(value)
	wd.AddDuration(duration)
	env.Store.Save(wd)
	return cli.OK
}
