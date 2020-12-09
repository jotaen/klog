package commands

import (
	"fmt"
	"klog/cli/lib"
)

func List(env lib.Environment, args []string) int {
	list, _ := env.Store.List()
	for _, date := range list {
		fmt.Printf("%v\n", date.ToString())
	}
	return lib.OK
}
