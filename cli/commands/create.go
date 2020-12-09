package commands

import (
	"klog/cli/lib"
	"klog/datetime"
	"klog/workday"
	"time"
)

func Create(env lib.Environment, args []string) int {
	today, _ := datetime.CreateDateFromTime(time.Now())
	wd := workday.Create(today)
	env.Store.Save(wd)
	return lib.OK
}
