package cli

import (
	"klog/app"
	"klog/project"
)

type Command struct {
	Main        func(app.Environment, project.Project, []string) int
	Name        string
	Alias       []string
	Description string
}
