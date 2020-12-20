package app

import (
	"klog/project"
	"klog/workday"
	"os/exec"
)

func OpenInEditor(project project.Project, workDay workday.WorkDay) error {
	props := project.GetFileProps(workDay)
	// open -t ...
	cmd := exec.Command("subl", props.Path)
	return cmd.Run()
}

func OpenInFileBrowser(project project.Project) error {
	cmd := exec.Command("open", project.Path())
	return cmd.Run()
}
