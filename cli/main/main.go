package main

import (
	"klog/cli"
	"os"
)

func main() {
	path, _ := os.Getwd()
	cli.Execute(path + "/tmp/cli", os.Args[1:])
}
