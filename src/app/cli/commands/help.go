package commands

import (
	"fmt"
	"klog/app"
	"klog/app/cli"
)

var Help cli.Command

func init() {
	Help = cli.Command{
		Name:        "help",
		Description: "Display help",
		Main:        help,
	}
}

func help(_ app.Service, _ []string) int {
	fmt.Printf("Help!\n")
	return cli.OK
}
