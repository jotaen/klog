package withenv

import (
	"klog/app/cli"
	"klog/project"
	"os"
)

func WithEnvironment(fn func(environment cli.Environment)) {
	path := "./tmp/test"
	os.RemoveAll(path)
	os.MkdirAll(path, os.ModePerm)
	st, err := project.NewProject(path)
	if err != nil {
		panic("Could not create project")
	}
	env := cli.Environment{
		WorkDir: path,
		Store:   st,
	}
	fn(env)
	os.RemoveAll(path)
	os.Remove("./tmp")
}
