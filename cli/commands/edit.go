package commands

import (
	"fmt"
	"klog/cli/lib"
	"klog/datetime"
	"klog/store"
	"os"
	"os/exec"
	"time"
)

func Edit(env lib.Environment, args []string) int {
	today, _ := datetime.CreateDateFromTime(time.Now())
	wd, err := env.Store.Get(today)
	if err != nil {
		fmt.Println("No no no no no no!")
		return lib.FILE_NOT_FOUND
	}

	file := env.Store.GetFileProps(wd)
	openEditor(file)
	return lib.OK
}

func openEditor(file store.FileProps) {
	cmd := exec.Command("vi", file.Path)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	err := cmd.Run()
	if err != nil {
		fmt.Println("Could not open editor")
	}
}
