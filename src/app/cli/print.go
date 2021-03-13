package cli

import (
	"klog/app"
	"klog/app/cli/lib"
	"klog/parser"
)

type Print struct {
	lib.FilterArgs
	lib.SortArgs
	lib.WarnArgs
	lib.InputFilesArgs
}

func (opt *Print) Run(ctx app.Context) error {
	records, err := ctx.ReadInputs(opt.File...)
	if err != nil {
		return err
	}
	if len(records) == 0 {
		return nil
	}
	now := ctx.Now()
	records = opt.ApplyFilter(now, records)
	records = opt.ApplySort(records)
	ctx.Print("\n" + parser.SerialiseRecords(&lib.Styler, records...) + "\n")

	ctx.Print(opt.WarnArgs.ToString(now, records))
	return nil
}
