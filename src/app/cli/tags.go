package cli

import (
	"github.com/jotaen/klog/src/app"
	"github.com/jotaen/klog/src/app/cli/lib"
	"github.com/jotaen/klog/src/app/cli/lib/terminalformat"
	"github.com/jotaen/klog/src/service"
)

type Tags struct {
	lib.FilterArgs
	lib.WarnArgs
	lib.NoStyleArgs
	lib.InputFilesArgs
}

func (opt *Tags) Run(ctx app.Context) error {
	opt.NoStyleArgs.Apply(&ctx)
	records, err := ctx.ReadInputs(opt.File...)
	if err != nil {
		return err
	}
	now := ctx.Now()
	records = opt.ApplyFilter(now, records)
	totalByTag := service.AggregateTotalsByTags(records...)
	if len(totalByTag) == 0 {
		return nil
	}
	table := terminalformat.NewTable(2, " ")
	for _, t := range totalByTag {
		if t.Tag.Value() == "" {
			table.CellL("#" + t.Tag.Name())
		} else {
			table.CellL(" " + t.Tag.Value())
		}
		table.CellL(ctx.Serialiser().Duration(t.Total))
	}
	table.Collect(ctx.Print)
	opt.WarnArgs.PrintWarnings(ctx, records)
	return nil
}
