package main

import (
	"klog/cli/core"
	"os"
)

func main() {
	path, _ := os.Getwd()
	code := core.Execute(path+"/tmp", os.Args[1:])
	os.Exit(code)
}
