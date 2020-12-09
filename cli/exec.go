package cli

import (
	"fmt"
	"klog/cli/commands"
	"klog/cli/lib"
	klogstore "klog/store"
)

type cmd func(klogstore.Store) int

var cmdDict map[string]cmd = map[string]cmd{
	"list":   commands.List,
	"create": commands.Create,
	"new":    commands.Create,
	"edit":   commands.Edit,
	"open":   commands.Edit,
	"start":  commands.Start,
}

type Environment struct {
	WorkDir string
}

func Execute(env Environment, args []string) int {
	store, err := klogstore.CreateFsStore(env.WorkDir)
	if err != nil {
		fmt.Printf("Project not found")
		return lib.PROJECT_PATH_INVALID
	}
	c := cmdDict[args[0]]
	return c(store)
}
