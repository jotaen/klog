package main

import (
	"klog/app/cli"
	"os"
)

func main() {
	code := cli.Execute()
	os.Exit(code)
}
