package main

import (
	"fmt"
	"os"
	klogstore "klog/store"
)

func main() {
	path, _ := os.Getwd()
	if len(os.Args) >= 2 {
		path += "/" + os.Args[1]
	}
	store, _ := klogstore.CreateFsStore(path)
	list, _ := store.List()
	for _, date := range list {
		fmt.Printf("%v\n", date.ToString())
	}
}
