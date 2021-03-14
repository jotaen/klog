package cli

import (
	. "klog"
	"klog/app"
	"klog/app/cli/lib"
	"klog/parser"
)

type Track struct {
	lib.AtDateArgs
	Entry string `arg required help:"A time entry, optionally with summary (might require quoting)"`
	lib.OutputFileArgs
}

func (opt *Track) Run(ctx app.Context) error {
	targetFile, err := opt.OutputFile(ctx)
	if err != nil {
		return err
	}
	pr, err := ctx.ReadFileInput(targetFile)
	if err != nil {
		return err
	}
	date := opt.AtDate()
	record, contents, err := pr.AddEntry(
		"No record at date "+date.ToString(),
		func(r Record) bool { return r.Date().IsEqualTo(date) },
		func(r Record) string { return opt.Entry })
	if err != nil {
		return err
	}
	err = ctx.WriteFile(targetFile, contents)
	if err != nil {
		return err
	}
	ctx.Print(parser.SerialiseRecords(&lib.Styler, record))
	return nil
}
