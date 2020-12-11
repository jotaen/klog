package commands

import (
	"fmt"
	"klog/cli"
)

var List cli.Command

func init() {
	List = cli.Command{
		Name:        "list",
		Alias:       []string{},
		Description: "List all entries in this project",
		Main:        list,
	}
}

func list(env cli.Environment, args []string) int {
	list, _ := env.Store.List()
	for _, date := range list {
		fmt.Printf("%v\n", date.ToString())
	}
	return cli.OK
}
