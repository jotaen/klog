package main

import (
	main2 "klog/app/cli/exec"
	"os"
)

func main() {
	path, _ := os.Getwd()
	code := main2.Execute(path+"/tmp", os.Args[1:])
	os.Exit(code)
}
