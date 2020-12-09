package main

import (
	"klog/cli"
	"os"
)

func main() {
	path, _ := os.Getwd()
	env := cli.Environment{
		WorkDir: path + "/tmp/cli",
	}
	if len(os.Args) > 1 {
		cli.Execute(env, os.Args[1:])
	}
}
