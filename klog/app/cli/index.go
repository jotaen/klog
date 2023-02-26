/*
Package cli contains handlers for all available commands.
*/
package cli

import (
	"github.com/jotaen/klog/klog/app"
	kc "github.com/jotaen/kong-completion"
)

type Cli struct {
	Default Default `hidden:"" cmd:"" default:"withargs" help:""`

	// Evaluate Files
	Print  Print  `cmd:"" name:"print" group:"Evaluate Files" help:"Pretty-prints records"`
	Total  Total  `cmd:"" name:"total" group:"Evaluate Files" help:"Evaluates the total time"`
	Report Report `cmd:"" name:"report" group:"Evaluate Files" help:"Prints an aggregated calendar report"`
	Tags   Tags   `cmd:"" name:"tags" group:"Evaluate Files" help:"Prints total times aggregated by tags"`
	Today  Today  `cmd:"" name:"today" group:"Evaluate Files" help:"Evaluates the current day"`

	// Manipulate Files
	Track  Track  `cmd:"" name:"track" group:"Manipulate Files" help:"Adds a new entry to a record"`
	Start  Start  `cmd:"" name:"start" group:"Manipulate Files" aliases:"in" help:"Starts a new open time range"`
	Stop   Stop   `cmd:"" name:"stop" group:"Manipulate Files" aliases:"out" help:"Closes the open time range"`
	Pause  Pause  `cmd:"" name:"pause" group:"Manipulate Files" help:"Pauses the open time range"`
	Create Create `cmd:"" name:"create" group:"Manipulate Files" help:"Creates a new, empty record"`

	// Manage Files
	Bookmarks Bookmarks `cmd:"" name:"bookmarks" group:"Manage Files" aliases:"bk" help:"Named aliases for often-used files"`
	Bookmark  Bookmarks `cmd:"" name:"bookmark" hidden:"" help:"(Alias)"` // Hidden alias for convenience / typo
	Edit      Edit      `cmd:"" name:"edit" group:"Manage Files" help:"Opens a file or bookmark in your editor"`
	Goto      Goto      `cmd:"" name:"goto" group:"Manage Files" help:"Opens the file explorer at a file or bookmark"`

	// Misc
	Version    Version       `cmd:"" name:"version" group:"Misc" help:"Prints version info and check for updates"`
	Config     Config        `cmd:"" name:"config" group:"Misc" help:"Prints the current configuration"`
	Info       Info          `cmd:"" name:"info" group:"Misc" help:"Prints information about klog"`
	Json       Json          `cmd:"" name:"json" group:"Misc" help:"Converts records to JSON"`
	Completion kc.Completion `cmd:"" name:"completion" group:"Misc" help:"Outputs shell code for enabling tab completion"`
}

const DESCRIPTION = "klog: command line app for time tracking with plain-text files.\n" +
	"Run with --help to learn usage.\n" +
	"Documentation online at " + KLOG_WEBSITE_URL

type Default struct {
	Version bool `short:"v" name:"version" help:"Alias for 'klog version'"`
}

func (opt *Default) Run(ctx app.Context) app.Error {
	if opt.Version {
		versionCmd := Version{}
		return versionCmd.Run(ctx)
	}
	ctx.Print(DESCRIPTION + "\n")
	return nil
}
