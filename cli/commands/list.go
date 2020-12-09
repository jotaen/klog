package commands

import (
	"fmt"
	"klog/cli/lib"
	"klog/store"
)

func List(st store.Store) int {
	list, _ := st.List()
	for _, date := range list {
		fmt.Printf("%v\n", date.ToString())
	}
	return lib.OK
}
