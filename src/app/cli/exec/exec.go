// Must be in separate package due to import cycle
package exec

import (
	"klog/app"
	"klog/app/cli"
	"klog/app/cli/commands"
)

var allCommands = []cli.Command{
	commands.Help,
	commands.Edit,
	commands.Start,
	commands.Widget,
}

func Execute(args []string) int {
	service, err := app.NewServiceWithConfigFiles()
	if err != nil {
		return cli.INITIALISATION_ERROR
	}
	if len(args) == 0 {
		args = []string{"help"}
	}
	reqSubCmd := args[0]
	for _, cmd := range allCommands {
		if cmd.Name == reqSubCmd {
			return cmd.Main(service, args)
		}
	}
	return cli.SUBCOMMAND_NOT_FOUND
}
