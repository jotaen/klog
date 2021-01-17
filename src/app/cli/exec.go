// Must be in separate package due to import cycle
package cli

import (
	"klog/app"
)

type Command struct {
	Main        func(app.Context, []string) int
	Name        string
	Description string
}

const (
	OK                   = 0
	EXECUTION_FAILED     = 1
	INITIALISATION_ERROR = 2
	SUBCOMMAND_NOT_FOUND = 3
	INVALID_CLI_ARGS     = 4
)

var allCommands = []Command{
	Help,
	Print,
	Edit,
	Start,
	Widget,
}

func Execute(args []string) int {
	ctx, err := app.NewContextFromEnv()
	if err != nil {
		return INITIALISATION_ERROR
	}
	if len(args) == 0 {
		args = []string{"help"}
	}
	subcommand := args[0]
	for _, cmd := range allCommands {
		if cmd.Name == subcommand {
			return cmd.Main(*ctx, args[1:])
		}
	}
	return SUBCOMMAND_NOT_FOUND
}
