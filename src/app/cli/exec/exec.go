package exec

import (
	"fmt"
	"klog/app"
	"klog/app/cli"
	"klog/app/cli/commands"
	"klog/project"
)

var allCommands = []cli.Command{
	commands.Create,
	commands.Edit,
	commands.List,
	commands.Log,
	commands.Start,
	commands.Tray,
}

func Execute(workDir string, args []string) int {
	st, err := project.NewProject(workDir)
	if err != nil {
		fmt.Printf("Project not found")
		return cli.PROJECT_PATH_INVALID
	}
	env := app.NewEnvironment("~")
	reqSubCmd := args[0]
	for _, cmd := range allCommands {
		if cmd.Name == reqSubCmd {
			return cmd.Main(env, st, args)
		}
	}
	return cli.SUBCOMMAND_NOT_FOUND
}
