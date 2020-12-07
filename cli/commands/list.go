package commands

import (
	"fmt"
	"klog/store"
)

func List(st store.Store) {
	list, _ := st.List()
	for _, date := range list {
		fmt.Printf("%v\n", date.ToString())
	}
}
