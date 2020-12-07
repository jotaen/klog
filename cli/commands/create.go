package commands

import (
	"klog/datetime"
	"klog/store"
	"klog/workday"
	"time"
)

func Create(st store.Store) {
	now := time.Now()
	date, _ := datetime.CreateDate(now.Year(), int(now.Month()), now.Day())
	wd := workday.Create(date)
	st.Save(wd)
}
