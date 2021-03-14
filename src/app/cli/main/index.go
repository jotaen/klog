package main

import . "klog/app/cli"

type cli struct {
	// Evaluate
	Print  Print  `cmd group:"Evaluate" help:"Pretty-print records"`
	Total  Total  `cmd group:"Evaluate" help:"Evaluate the total time"`
	Report Report `cmd group:"Evaluate" help:"Print a calendar report summarising all days"`
	Tags   Tags   `cmd group:"Evaluate" help:"Print total times aggregated by tags"`
	Now    Now    `cmd group:"Evaluate" help:"Evaluate today’s record (including potential open ranges)"`

	// Manipulate
	Append Append `cmd group:"Manipulate" hidden help:"Appends a new record to a file (based on templates)"`
	Track  Track  `cmd group:"Manipulate" hidden help:"Add a new entry to a record"`
	Start  Start  `cmd group:"Manipulate" hidden aliases:"in" help:"Start open time range"`

	// Misc
	Bookmark Bookmark `cmd group:"Misc" help:"Default file that klog reads from"`
	Json     Json     `cmd group:"Misc" help:"Convert records to JSON"`
	Widget   Widget   `cmd group:"Misc" help:"Start menu bar widget (MacOS only)"`
	Version  Version  `cmd group:"Misc" help:"Print version info and check for updates"`
}
