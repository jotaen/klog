package withenv

import (
	"klog/app"
	"klog/project"
	"os"
)

func WithEnvironment(fn func(app.Environment, project.Project)) {
	path := "./tmp/test"
	os.RemoveAll(path)
	os.MkdirAll(path, os.ModePerm)
	p, err := project.NewProject(path)
	if err != nil {
		panic("Could not create project")
	}
	env := app.NewEnvironment("~")
	fn(env, p)
	os.RemoveAll(path)
	os.Remove("./tmp")
}
