package main

import (
	"klog/app/cli/exec"
	"os"
)

func main() {
	code := exec.Execute(os.Args[1:])
	os.Exit(code)
}
