/*
Package cli contains handlers for all available commands.
*/
package cli

import kcmpl "github.com/jotaen/kong-completion"

type Cli struct {
	// Evaluate Files
	Print  Print  `cmd:"" group:"Evaluate Files" help:"Pretty-prints records"`
	Total  Total  `cmd:"" group:"Evaluate Files" help:"Evaluates the total time"`
	Report Report `cmd:"" group:"Evaluate Files" help:"Prints a calendar report summarising all days"`
	Tags   Tags   `cmd:"" group:"Evaluate Files" help:"Prints total times aggregated by tags"`
	Today  Today  `cmd:"" group:"Evaluate Files" help:"Evaluates the current day"`

	// Manipulate Files
	Track  Track  `cmd:"" group:"Manipulate Files" help:"Adds a new entry to a record"`
	Start  Start  `cmd:"" group:"Manipulate Files" aliases:"in" help:"Starts a new open time range"`
	Stop   Stop   `cmd:"" group:"Manipulate Files" aliases:"out" help:"Closes the open time range"`
	Pause  Pause  `cmd:"" group:"Manipulate Files" help:"Pauses the open time range"`
	Create Create `cmd:"" group:"Manipulate Files" help:"Creates a new record"`

	// Manage Files
	Bookmarks Bookmarks `cmd:"" group:"Manage Files" aliases:"bk" help:"Named aliases for often-used files"`
	Bookmark  Bookmarks `cmd:"" hidden:"" help:"(Alias)"` // Hidden alias for convenience / typo
	Edit      Edit      `cmd:"" group:"Manage Files" help:"Opens a file or bookmark in your editor"`
	Goto      Goto      `cmd:"" group:"Manage Files" help:"Opens the file explorer at a file or bookmark"`

	// Misc
	Version    Version          `cmd:"" group:"Misc" help:"Prints version info and check for updates"`
	Info       Info             `cmd:"" group:"Misc" default:"withargs" help:"Displays meta info about klog"`
	Config     Config           `cmd:"" group:"Misc" help:"Prints the current configuration"`
	Json       Json             `cmd:"" group:"Misc" help:"Converts records to JSON"`
	Completion kcmpl.Completion `cmd:"" group:"Misc" help:"Outputs shell code for enabling tab completion"`
}
