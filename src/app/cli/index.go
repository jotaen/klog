package cli

type Cli struct {
	// Evaluate
	Print  Print  `cmd group:"Evaluate" help:"Pretty-prints records"`
	Total  Total  `cmd group:"Evaluate" help:"Evaluates the total time"`
	Report Report `cmd group:"Evaluate" help:"Prints a calendar report summarising all days"`
	Tags   Tags   `cmd group:"Evaluate" help:"Prints total times aggregated by tags"`
	Today  Today  `cmd group:"Evaluate" help:"Evaluates the current day"`

	// Manipulate
	Track  Track  `cmd group:"Manipulate" help:"Adds a new entry to a record"`
	Start  Start  `cmd group:"Manipulate" aliases:"in" help:"Starts open time range"`
	Stop   Stop   `cmd group:"Manipulate" aliases:"out" help:"Closes open time range"`
	Create Create `cmd group:"Manipulate" help:"Creates a new record"`

	// Bookmarks
	Bookmarks Bookmarks `cmd group:"Bookmarks" help:"Named aliases for often-used files"`
	Bookmark  Bookmarks `cmd group:"Bookmarks" hidden help:"Alias"`
	Bm        Bookmarks `cmd group:"Bookmarks" hidden help:"Alias"`

	// Misc
	Edit    Edit    `cmd group:"Misc" help:"Opens a file or bookmark in your editor"`
	Json    Json    `cmd group:"Misc" help:"Converts records to JSON"`
	Widget  Widget  `cmd group:"Misc" help:"Starts menu bar widget (MacOS only)"`
	Version Version `cmd group:"Misc" help:"Prints version info and check for updates"`

	// Default command for displaying info text (hidden)
	Info Info `cmd default:"withargs" hidden:"1"`
}
