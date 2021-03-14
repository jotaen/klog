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
	return handleAddEntry(
		opt.OutputFileArgs,
		ctx,
		func(pr *parser.ParseResult) (Record, string, error) {
			date := opt.AtDate(ctx.Now())
			return pr.AddEntry(
				"No record at date "+date.ToString(),
				func(r Record) bool { return r.Date().IsEqualTo(date) },
				func(r Record) string { return opt.Entry },
			)
		},
	)
}

func handleAddEntry(
	fileArgs lib.OutputFileArgs,
	ctx app.Context,
	handler func(*parser.ParseResult) (Record, string, error),
) error {
	targetFile, err := fileArgs.OutputFile(ctx)
	if err != nil {
		return err
	}
	pr, err := ctx.ReadFileInput(targetFile)
	if err != nil {
		return err
	}
	record, contents, err := handler(pr)
	if err != nil {
		return err
	}
	err = ctx.WriteFile(targetFile, contents)
	if err != nil {
		return err
	}
	ctx.Print("\n" + parser.SerialiseRecords(&lib.Styler, record) + "\n")
	return nil
}
