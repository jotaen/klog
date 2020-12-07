package main

import (
	"klog/cli"
	klogstore "klog/store"
	"os"
)

func main() {
	path, _ := os.Getwd()
	store, err := klogstore.CreateFsStore(path + "/tmp/cli")
	if err != nil {
		os.Exit(1)
	}
	if len(os.Args) > 1 {
		cli.Exec(store, os.Args[1])
	}
}
