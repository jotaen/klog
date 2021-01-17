package cli

import (
	"fmt"
	"klog/app"
	"klog/parser"
)

var Print = Command{
	Name:        "print",
	Description: "Print a file",
	Main:        print,
}

func print(ctx app.Context, args []string) int {
	if len(args) == 0 {
		fmt.Println("Please specify a file")
		return INVALID_CLI_ARGS
	}
	rs, err := ctx.Read("../" + args[0])
	if err != nil {
		fmt.Println(err)
		return EXECUTION_FAILED
	}
	fmt.Println(parser.SerialiseRecords(rs))
	return OK
}
