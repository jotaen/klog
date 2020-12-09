package commands

import (
	"klog/cli/lib"
	"klog/datetime"
	"klog/store"
	"klog/workday"
	"time"
)

func Create(st store.Store) int {
	now := time.Now()
	today, _ := datetime.CreateDate(now.Year(), int(now.Month()), now.Day())
	wd := workday.Create(today)
	st.Save(wd)
	return lib.OK
}
