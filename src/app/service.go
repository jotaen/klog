package app

import (
	"klog/datetime"
	"klog/project"
	"klog/workday"
	"os/exec"
	"time"
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

func Start(project project.Project, start time.Time) (workday.WorkDay, error) {
	today, _ := datetime.NewDateFromTime(start)
	wd, _ := project.Get(today)
	if wd == nil {
		wd = workday.NewWorkDay(today)
	}
	startTime, _ := datetime.CreateTimeFromTime(start)
	wd.StartOpenRange(startTime)
	project.Save(wd)
	return wd, nil
}

func Stop(project project.Project, end time.Time) (workday.WorkDay, error) {
	today, _ := datetime.NewDateFromTime(end)
	wd, err := project.Get(today)
	if wd == nil {
		return nil, err[0] // todo
	}
	startTime, _ := datetime.CreateTimeFromTime(end)
	wd.EndOpenRange(startTime)
	project.Save(wd)
	return wd, nil
}
