package exec

import (
	"klog/app"
	"klog/app/cli"
	"klog/app/cli/commands"
)

var allCommands = []cli.Command{
	commands.Create,
	commands.Edit,
	commands.Start,
	commands.Widget,
}

func Execute(args []string) int {
	service := app.NewService(nil) // TODO
	reqSubCmd := args[0]
	for _, cmd := range allCommands {
		if cmd.Name == reqSubCmd {
			return cmd.Main(service, args)
		}
	}
	return cli.SUBCOMMAND_NOT_FOUND
}
