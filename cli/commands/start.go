package commands

import (
	"klog/cli/lib"
	"klog/datetime"
	"klog/store"
	"time"
)

func Start(st store.Store) int {
	now := time.Now()
	today, _ := datetime.CreateDateFromTime(now)
	wd, err := st.Get(today)
	if err != nil {
		// todo create new
		return 182763
	}
	nowTime, _ := datetime.CreateTimeFromTime(now)
	wd.AddOpenRange(nowTime)
	st.Save(wd)
	return lib.OK
}
