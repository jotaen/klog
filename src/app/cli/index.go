/*
Package cli contains handlers for all available commands.
*/
package cli

type Cli struct {
	// Evaluate
	Print  Print  `cmd:"" group:"Evaluate" help:"Pretty-prints records"`
	Total  Total  `cmd:"" group:"Evaluate" help:"Evaluates the total time"`
	Report Report `cmd:"" group:"Evaluate" help:"Prints a calendar report summarising all days"`
	Tags   Tags   `cmd:"" group:"Evaluate" help:"Prints total times aggregated by tags"`
	Today  Today  `cmd:"" group:"Evaluate" help:"Evaluates the current day"`

	// Manipulate
	Track  Track  `cmd:"" group:"Manipulate" help:"Adds a new entry to a record"`
	Start  Start  `cmd:"" group:"Manipulate" aliases:"in" help:"Starts a new open time range"`
	Stop   Stop   `cmd:"" group:"Manipulate" aliases:"out" help:"Closes the open time range"`
	Pause  Pause  `cmd:"" group:"Manipulate" help:"Pauses the open time range"`
	Create Create `cmd:"" group:"Manipulate" help:"Creates a new record"`

	// Bookmarks
	Bookmarks Bookmarks `cmd:"" group:"Bookmarks (bk)" help:"Named aliases for often-used files"`
	Bookmark  Bookmarks `cmd:"" group:"Bookmarks" hidden:"" help:"Alias"`
	Bm        Bookmarks `cmd:"" group:"Bookmarks" hidden:"" help:"Alias"`
	Bk        Bookmarks `cmd:"" group:"Bookmarks" hidden:"" help:"Alias"`

	// Misc
	Edit       Edit       `cmd:"" group:"Misc" help:"Opens a file or bookmark in your editor"`
	Goto       Goto       `cmd:"" group:"Misc" help:"Opens the file explorer at the given location"`
	Json       Json       `cmd:"" group:"Misc" help:"Converts records to JSON"`
	Info       Info       `cmd:"" group:"Misc" default:"withargs" help:"Displays meta info about klog"`
	Version    Version    `cmd:"" group:"Misc" help:"Prints version info and check for updates"`
	Completion Completion `cmd:"" group:"Misc" help:"Output shell code for initialising shell completions for klog"`
}
