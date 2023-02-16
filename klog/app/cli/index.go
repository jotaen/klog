/*
Package cli contains handlers for all available commands.
*/
package cli

import kcmpl "github.com/jotaen/kong-completion"

type Cli struct {
	// Evaluate
	Print  Print  `cmd:"" group:"Evaluate Files" help:"Pretty-prints records"`
	Total  Total  `cmd:"" group:"Evaluate Files" help:"Evaluates the total time"`
	Report Report `cmd:"" group:"Evaluate Files" help:"Prints a calendar report summarising all days"`
	Tags   Tags   `cmd:"" group:"Evaluate Files" help:"Prints total times aggregated by tags"`
	Today  Today  `cmd:"" group:"Evaluate Files" help:"Evaluates the current day"`

	// Manipulate
	Track  Track  `cmd:"" group:"Manipulate Files" help:"Adds a new entry to a record"`
	Start  Start  `cmd:"" group:"Manipulate Files" aliases:"in" help:"Starts a new open time range"`
	Stop   Stop   `cmd:"" group:"Manipulate Files" aliases:"out" help:"Closes the open time range"`
	Pause  Pause  `cmd:"" group:"Manipulate Files" help:"Pauses the open time range"`
	Create Create `cmd:"" group:"Manipulate Files" help:"Creates a new record"`

	// Bookmarks
	Bookmarks Bookmarks `cmd:"" group:"Bookmarks (bk)" help:"Named aliases for often-used files"`
	Bookmark  Bookmarks `cmd:"" group:"Bookmarks" hidden:"" help:"Alias"`
	Bm        Bookmarks `cmd:"" group:"Bookmarks" hidden:"" help:"Alias"`
	Bk        Bookmarks `cmd:"" group:"Bookmarks" hidden:"" help:"Alias"`

	// Misc
	Edit       Edit             `cmd:"" group:"Misc" help:"Opens a file or bookmark in your editor"`
	Goto       Goto             `cmd:"" group:"Misc" help:"Opens the file explorer at the given location"`
	Config     Config           `cmd:"" group:"Misc" help:"Prints the current configuration"`
	Json       Json             `cmd:"" group:"Misc" help:"Converts records to JSON"`
	Info       Info             `cmd:"" group:"Misc" default:"withargs" help:"Displays meta info about klog"`
	Version    Version          `cmd:"" group:"Misc" help:"Prints version info and check for updates"`
	Completion kcmpl.Completion `cmd:"" group:"Misc" help:"Outputs shell code for enabling tab completion"`
}
