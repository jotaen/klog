package commands

import (
	"fmt"
	"klog/cli/lib"
	"klog/datetime"
	klogstore "klog/store"
	"os"
	"os/exec"
	"time"
)

func Edit(store klogstore.Store) int {
	today, _ := datetime.CreateDateFromTime(time.Now())
	wd, err := store.Get(today)
	if err != nil {
		fmt.Println("No no no no no no!")
		return lib.FILE_NOT_FOUND
	}

	file := store.GetFileProps(wd)
	openEditor(file)
	return lib.OK
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
