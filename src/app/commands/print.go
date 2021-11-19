package commands

import (
	"github.com/jotaen/klog/src/app"
	"github.com/jotaen/klog/src/app/lib"
)

type Print struct {
	lib.FilterArgs
	lib.SortArgs
	lib.WarnArgs
	lib.NoStyleArgs
	lib.InputFilesArgs
}

func (opt *Print) Help() string {
	return `The output is syntax-highlighted and the formatting is slightly sanitised.`
}

func (opt *Print) Run(ctx app.Context) error {
	opt.NoStyleArgs.Apply(&ctx)
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
	ctx.Print("\n" + ctx.Serialiser().SerialiseRecords(records...) + "\n")

	ctx.Print(opt.WarnArgs.ToString(now, records))
	return nil
}
