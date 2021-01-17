package main

import (
	"klog/app/cli"
	"os"
)

func main() {
	code := cli.Execute(os.Args[1:])
	os.Exit(code)
}
