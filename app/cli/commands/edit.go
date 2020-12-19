package commands

import (
	"fmt"
	"klog/app/cli"
	"klog/datetime"
	"klog/store"
	"os"
	"os/exec"
	"time"
)

var Edit cli.Command

func init() {
	Edit = cli.Command{
		Name:        "edit",
		Alias:       []string{"open"},
		Description: "Open entry in editor",
		Main:        edit,
	}
}

func edit(env cli.Environment, args []string) int {
	today, _ := datetime.NewDateFromTime(time.Now())
	wd, err := env.Store.Get(today)
	if err != nil {
		fmt.Println("No no no no no no!")
		return cli.FILE_NOT_FOUND
	}

	file := env.Store.GetFileProps(wd)
	openEditor(file)
	return cli.OK
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
