package cli

import (
	"github.com/jotaen/klog/klog/app"
	"github.com/jotaen/klog/klog/app/cli/lib"
	"github.com/jotaen/klog/klog/parser"
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
	now := ctx.Now()
	records = opt.ApplyFilter(now, records)
	if len(records) == 0 {
		return nil
	}
	records = opt.ApplySort(records)
	ctx.Print("\n" + parser.SerialiseRecords(ctx.Serialiser(), records...) + "\n")

	opt.WarnArgs.PrintWarnings(ctx, records)
	return nil
}
