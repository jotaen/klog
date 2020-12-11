package core

import (
	"fmt"
	"klog/cli"
	"klog/cli/commands"
	"klog/store"
)

type cmd func(cli.Environment, []string) int

var cmdDict = map[string]cmd{
	"list":   commands.List,
	"create": commands.Create,
	"new":    commands.Create,
	"edit":   commands.Edit,
	"open":   commands.Edit,
	"start":  commands.Start,
	"log":    commands.Log,
}

func Execute(workDir string, args []string) int {
	st, err := store.CreateFsStore(workDir)
	if err != nil {
		fmt.Printf("Project not found")
		return cli.PROJECT_PATH_INVALID
	}
	env := cli.Environment{
		WorkDir: workDir,
		Store:   st,
	}
	c := cmdDict[args[0]]
	return c(env, args)
}
