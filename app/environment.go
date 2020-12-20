package app

import (
	"klog/project"
)

type Environment interface {
	SavedProjects() []project.Project
}

type environment struct {
	projects []project.Project
}

func NewEnvironment(homeFolderPath string) Environment {
	// TODO
	defaultProject, _ := project.NewProject("./tmp")
	return environment{
		projects: []project.Project{defaultProject},
	}
}

func (c environment) SavedProjects() []project.Project {
	return c.projects
}
