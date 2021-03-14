package cli

import (
	. "klog"
	"klog/app"
	"klog/app/cli/lib"
	"klog/parser"
)

type Start struct {
	lib.AtDateArgs
	lib.OutputFileArgs
}

func (opt *Start) Run(ctx app.Context) error {
	return handleAddEntry(
		opt.OutputFileArgs,
		ctx,
		func(pr *parser.ParseResult) (Record, string, error) {
			date := opt.AtDate(ctx.Now())
			time := NewTimeFromTime(ctx.Now())
			return pr.AppendEntry(
				"No record at date "+date.ToString(),
				func(r Record) bool { return r.Date().IsEqualTo(date) },
				func(r Record) string { return time.ToString() + "-?" },
			)
		},
	)
}
