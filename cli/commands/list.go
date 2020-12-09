package commands

import (
	"fmt"
	"klog/cli"
)

func List(env cli.Environment, args []string) int {
	list, _ := env.Store.List()
	for _, date := range list {
		fmt.Printf("%v\n", date.ToString())
	}
	return cli.OK
}
