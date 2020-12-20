package commands

import (
	"fmt"
	"klog/app"
	"klog/app/cli"
	"klog/project"
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

func list(env app.Environment, project project.Project, args []string) int {
	list, _ := project.List()
	for _, date := range list {
		fmt.Printf("%v\n", date.ToString())
	}
	return cli.OK
}
