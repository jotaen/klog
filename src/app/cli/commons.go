package cli

type FileArgs struct {
	File string `arg optional name:"file" help:"File to read from"`
}

type FilterArgs struct {
	Tag    []string `short:"t" long:"tag" default:"\"TAG\"" help:"Only records that contain this tag"`
	Date   string   `short:"d" long:"date" default:"DATE" help:"Only records at this date"`
	After  string   `short:"a" long:"after" default:"DATE" help:"Only records at or after this date"`
	Before string   `short:"b" long:"before" default:"DATE" help:"Only records at or before this date"`
}
