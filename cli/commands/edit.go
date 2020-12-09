package commands

import (
	"fmt"
	"klog/datetime"
	klogstore "klog/store"
	"os"
	"os/exec"
	"time"
)

func Edit(store klogstore.Store) {
	now := time.Now()
	today, _ := datetime.CreateDate(now.Year(), int(now.Month()), now.Day())
	wd, err := store.Get(today)
	if err != nil {
		fmt.Println("No no no no no no!")
		return
	}

	file := store.GetFileProps(wd)
	openEditor(file)
}

func openEditor(file klogstore.FileProps) {
	cmd := exec.Command("vi", file.Path)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	err := cmd.Run()
	if err != nil {
		fmt.Println("Could not open editor")
	}
}
