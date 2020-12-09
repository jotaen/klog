package cli

import (
	"fmt"
	"klog/cli/commands"
	"klog/cli/lib"
	"klog/store"
)

type cmd func(lib.Environment, []string) int

var cmdDict map[string]cmd = map[string]cmd{
	"list":   commands.List,
	"create": commands.Create,
	"new":    commands.Create,
	"edit":   commands.Edit,
	"open":   commands.Edit,
	"start":  commands.Start,
	"log":    commands.Log,
}

func Execute(workDir string, args []string) int {
	store, err := store.CreateFsStore(workDir)
	if err != nil {
		fmt.Printf("Project not found")
		return lib.PROJECT_PATH_INVALID
	}
	env := lib.Environment{
		WorkDir: workDir,
		Store: store,
	}
	c := cmdDict[args[0]]
	return c(env, args[1:])
}
