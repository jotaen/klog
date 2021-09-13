package cli

import (
	"errors"
	"github.com/jotaen/klog/src/app"
)

type Cli struct {
	// Evaluate
	Print  Print  `cmd group:"Evaluate" help:"Pretty-prints records"`
	Total  Total  `cmd group:"Evaluate" help:"Evaluates the total time"`
	Report Report `cmd group:"Evaluate" help:"Prints a calendar report summarising all days"`
	Tags   Tags   `cmd group:"Evaluate" help:"Prints total times aggregated by tags"`
	Today  Today  `cmd group:"Evaluate" help:"Evaluate current day"`
	Now    Now    `cmd group:"Evaluate" hidden help:"Show overview of the current day"`

	// Manipulate
	Track  Track  `cmd group:"Manipulate" help:"Adds a new entry to a record"`
	Start  Start  `cmd group:"Manipulate" aliases:"in" help:"Starts open time range"`
	Stop   Stop   `cmd group:"Manipulate" aliases:"out" help:"Closes open time range"`
	Create Create `cmd group:"Manipulate" help:"Creates a new record"`

	// Misc
	Bookmark Bookmark `cmd group:"Misc" help:"Default file that klog reads from"`
	Json     Json     `cmd group:"Misc" help:"Converts records to JSON"`
	Widget   Widget   `cmd group:"Misc" help:"Starts menu bar widget (MacOS only)"`
	Version  Version  `cmd group:"Misc" help:"Prints version info and check for updates"`
}

// DEPRECATED
type Now struct {
	Today
}

func (opt *Now) Run(_ app.Context) error {
	return errors.New("`klog now` has been renamed to `klog today`. Note that the `--now` flag must be passed explicitly.")
}
