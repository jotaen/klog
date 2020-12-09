package commands

import (
	"klog/cli/lib"
	"klog/datetime"
	"klog/store"
	"klog/workday"
	"time"
)

func Create(st store.Store) int {
	today, _ := datetime.CreateDateFromTime(time.Now())
	wd := workday.Create(today)
	st.Save(wd)
	return lib.OK
}
