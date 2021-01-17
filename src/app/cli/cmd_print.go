package cli

import (
	"fmt"
	"klog/app"
)

var Print Command

func init() {
	Print = Command{
		Name:        "print",
		Description: "Print a file",
		Main:        print,
	}
}

func print(_ app.Context, args []string) int {
	if len(args) == 0 {
		fmt.Println("Please specify a file")
		return INVALID_CLI_ARGS
	}

	return OK
}
