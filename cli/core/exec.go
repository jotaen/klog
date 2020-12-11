package core

import (
	"fmt"
	"klog/cli"
	"klog/cli/commands"
	"klog/store"
)

var allCommands = []cli.Command{
	commands.Create,
	commands.Edit,
	commands.List,
	commands.Log,
	commands.Start,
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
	reqSubCmd := args[0]
	for _, cmd := range allCommands {
		if cmd.Name == reqSubCmd {
			return cmd.Main(env, args)
		}
	}
	return cli.SUBCOMMAND_NOT_FOUND
}
