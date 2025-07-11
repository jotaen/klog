/*
Package cli contains handlers for all available commands.
*/
package cli

import (
	"github.com/jotaen/klog/klog/app"
	kc "github.com/jotaen/kong-completion"
)

var INTRO_SUMMARY = "klog is a command-line tool for time tracking in a human-readable, plain-text file format.\nSee " + KLOG_WEBSITE_URL + " for documentation.\n"

// Guideline for help texts and descriptions:
// - Command and flag descriptions are phrased in imperative style, and they
//   end in a period. Examples:
//   - Pretty-print records.
//   - Sort output by date.
// - Code and literal values are wrapped in single quotes (')
// - Types of flag values are spelled in UPPER-CASE. The type is explained in
//   the flag description. For complex types, there should also be an example.

type Cli struct {
	Default Default `hidden:"" cmd:"" default:"withargs" help:""`

	// Evaluate Files
	Print  Print  `cmd:"" name:"print" group:"Evaluate Files" help:"Pretty-print records."`
	Total  Total  `cmd:"" name:"total" group:"Evaluate Files" help:"Evaluate the total time."`
	Report Report `cmd:"" name:"report" group:"Evaluate Files" help:"Print an aggregated calendar report."`
	Tags   Tags   `cmd:"" name:"tags" group:"Evaluate Files" help:"Print total times aggregated by tags."`
	Today  Today  `cmd:"" name:"today" group:"Evaluate Files" help:"Evaluate the current day."`

	// Manipulate Files
	Track  Track  `cmd:"" name:"track" group:"Manipulate Files" help:"Add a new entry to a record."`
	Start  Start  `cmd:"" name:"start" group:"Manipulate Files" aliases:"in" help:"Start a new open time range."`
	Stop   Stop   `cmd:"" name:"stop" group:"Manipulate Files" aliases:"out" help:"Close the open time range."`
	Pause  Pause  `cmd:"" name:"pause" group:"Manipulate Files" help:"Pause the open time range."`
	Switch Switch `cmd:"" name:"switch" group:"Manipulate Files" help:"Close open range and starts a new one."`
	Create Create `cmd:"" name:"create" group:"Manipulate Files" help:"Create a new, empty record."`

	// Manage Files
	Bookmarks Bookmarks `cmd:"" name:"bookmarks" group:"Manage Files" aliases:"bk" help:"Named aliases for often-used files."`
	Bookmark  Bookmarks `cmd:"" name:"bookmark" hidden:"" help:"(Alias)"` // Hidden alias for convenience / typo
	Edit      Edit      `cmd:"" name:"edit" group:"Manage Files" help:"Open a file or bookmark in your editor."`
	Goto      Goto      `cmd:"" name:"goto" group:"Manage Files" help:"Open the file explorer at a file or bookmark."`

	// Misc
	Version    Version       `cmd:"" name:"version" group:"Misc" help:"Print version info and check for updates."`
	Config     Config        `cmd:"" name:"config" group:"Misc" help:"Print the current configuration."`
	Info       Info          `cmd:"" name:"info" group:"Misc" help:"Print information about klog."`
	Json       Json          `cmd:"" name:"json" group:"Misc" help:"Convert records to JSON."`
	Completion kc.Completion `cmd:"" name:"completion" group:"Misc" help:"Output shell code for enabling tab completion."`
}

type Default struct {
	Version bool `short:"v" name:"version" help:"Alias for 'klog version'."`
}

func (opt *Default) Help() string {
	return INTRO_SUMMARY + `

Time-tracking data is stored in files ending in the '.klg' extension.
You can use the subcommands listed below to evaluate, manipulate and manage your files.
Use the '--help' flag on the subcommands to learn more.

You can specify input data in one of these 3 ways:
  - by passing the name of a file or a bookmark,
  - by piping data to stdin,
  - or by setting up a default bookmark.

Run 'klog bookmarks --help' to learn about bookmark usage.

Some general notes on flag usage:
  - For flags with values, you can either use a space or an equals sign as delimiter. E.g., both '--flag value' and '--flag=value' are fine.
  - For shorthand flags with values, you specify the value without a delimiter. E.g., '-n2' (if the long form is '--number 2').
  - For shorthand flags without values, you can compact them. E.g., '-abc' is the same as '-a -b -c'.
`
}

func (opt *Default) Run(ctx app.Context) app.Error {
	if opt.Version {
		versionCmd := Version{}
		return versionCmd.Run(ctx)
	}
	ctx.Print(INTRO_SUMMARY)
	ctx.Print("Run 'klog --help' to learn usage.\n")
	return nil
}
