package cli

import (
	. "klog"
	"klog/app"
	"klog/app/cli/lib"
	gotime "time"
)

type Track struct {
	Date  Date   `name:"date" group:"Filter" help:"The date at which to add the entry (defaults to today)"`
	Entry string `arg required help:"A time entry, optionally with summary (might require quoting)"`
	lib.OutputFileArgs
}

func (opt *Track) atDate() Date {
	if opt.Date != nil {
		return opt.Date
	}
	return NewDateFromTime(gotime.Now())
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
	today := opt.atDate()
	contents, err := pr.AddEntry(func(rs []Record) (int, string) {
		for i, r := range rs {
			if r.Date().IsEqualTo(today) {
				return i, opt.Entry
			}
		}
		return -1, ""
	})
	if err != nil {
		return err
	}
	return ctx.WriteFile(targetFile, contents)
}
