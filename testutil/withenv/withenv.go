package withenv

import (
	"klog/cli"
	"klog/store"
	"os"
)

func WithEnvironment(fn func(environment cli.Environment)) {
	path := "./tmp/test"
	os.RemoveAll(path)
	os.MkdirAll(path, os.ModePerm)
	st, err := store.NewFsStore(path)
	if err != nil {
		panic("Could not create store")
	}
	env := cli.Environment{
		WorkDir: path,
		Store:   st,
	}
	fn(env)
	os.RemoveAll(path)
	os.Remove("./tmp")
}
