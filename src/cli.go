package main

import (
	"klog/app/cli/exec"
	"os"
)

func main() {
	path, _ := os.Getwd()
	code := exec.Execute(path+"/tmp", os.Args[1:])
	os.Exit(code)
}
