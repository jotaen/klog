package app

import "klog/project"

type Configuration interface {
	SavedProjects() []project.Project
}

type configuration struct {
	projects []project.Project
}

func NewConfiguration(homeFolderPath string) Configuration {
	// TODO
	defaultProject, _ := project.NewProject("./tmp")
	return configuration{
		projects: []project.Project{defaultProject},
	}
}

func (c configuration) SavedProjects() []project.Project {
	return c.projects
}
