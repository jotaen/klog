package cli

import (
	"fmt"
	"klog/app"
)

var Help Command

func init() {
	Help = Command{
		Name:        "help",
		Description: "Display help",
		Main:        help,
	}
}

func help(_ app.Context, _ []string) int {
	fmt.Printf("Help!\n")
	return OK
}
