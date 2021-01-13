package commands

import (
	"fmt"
	"klog/app"
	"klog/app/cli"
	"klog/datetime"
	"klog/project"
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

func edit(env app.Environment, project project.Project, args []string) int {
	today, _ := datetime.NewDateFromTime(time.Now())
	wd, err := project.Get(today)
	if err != nil {
		fmt.Println("No no no no no no!")
		return cli.FILE_NOT_FOUND
	}

	file := project.GetFileProps(wd)
	openEditor(file)
	return cli.OK
}

func openEditor(file project.FileProps) {
	cmd := exec.Command("vi", file.Path)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	err := cmd.Run()
	if err != nil {
		fmt.Println("Could not open editor")
	}
}
