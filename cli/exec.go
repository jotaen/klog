package cli

import (
	"fmt"
	"klog/cli/commands"
	klogstore "klog/store"
)

type cmd func(klogstore.Store)

var cmdDict map[string]cmd = map[string]cmd{
	"list":   commands.List,
	"create": commands.Create,
	"new":    commands.Create,
	"edit":   commands.Edit,
	"open":   commands.Edit,
}

type Environment struct {
	WorkDir string
}

func Execute(env Environment, args []string) int {
	store, err := klogstore.CreateFsStore(env.WorkDir)
	if err != nil {
		fmt.Printf("Project not found")
		return 1
	}
	c := cmdDict[args[0]]
	if c == nil {
		fmt.Printf("Project not found")
		return 1
	}
	c(store)
	return 0
}
