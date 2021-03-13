package cli

import (
	"klog/app"
	"klog/app/cli/lib"
	"klog/parser"
)

type Print struct {
	FilterArgs
	SortArgs
	WarnArgs
	InputFilesArgs
}

func (opt *Print) Run(ctx app.Context) error {
	records, err := ctx.RetrieveRecords(opt.File...)
	if err != nil {
		return err
	}
	if len(records) == 0 {
		return nil
	}
	now := ctx.Now()
	records = opt.filter(now, records)
	records = opt.sort(records)
	ctx.Print("\n" + parser.SerialiseRecords(&lib.Styler, records...) + "\n")

	ctx.Print(opt.WarnArgs.ToString(now, records))
	return nil
}
