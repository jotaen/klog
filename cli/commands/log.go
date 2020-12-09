package commands

import (
	"klog/cli/lib"
	"klog/datetime"
	"time"
)

func Log(env lib.Environment, args []string) int {
	now := time.Now()
	today, _ := datetime.CreateDateFromTime(now)
	wd, err := env.Store.Get(today)
	if err != nil {
		// todo create new
		return 182763
	}
	duration, _ := datetime.CreateDurationFromString(args[0])
	wd.AddTime(duration)
	env.Store.Save(wd)
	return lib.OK
}
