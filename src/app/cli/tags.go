package cli

import (
	"github.com/jotaen/klog/src/app"
	"github.com/jotaen/klog/src/app/cli/lib"
	"github.com/jotaen/klog/src/app/cli/lib/terminalformat"
	"github.com/jotaen/klog/src/service"
)

type Tags struct {
	Values bool `name:"values" short:"v" help:"Display breakdown of tag values"`
	lib.FilterArgs
	lib.WarnArgs
	lib.NowArgs
	lib.DecimalArgs
	lib.NoStyleArgs
	lib.InputFilesArgs
}

func (opt *Tags) Run(ctx app.Context) error {
	opt.DecimalArgs.Apply(&ctx)
	opt.NoStyleArgs.Apply(&ctx)
	records, err := ctx.ReadInputs(opt.File...)
	if err != nil {
		return err
	}
	now := ctx.Now()
	records = opt.ApplyFilter(now, records)
	records, nErr := opt.ApplyNow(now, records...)
	if nErr != nil {
		return nErr
	}
	totalByTag := service.AggregateTotalsByTags(records...)
	if len(totalByTag) == 0 {
		return nil
	}
	numberOfColumns := 2
	if opt.Values {
		numberOfColumns = 3
	}
	table := terminalformat.NewTable(numberOfColumns, " ")
	for _, t := range totalByTag {
		totalString := ctx.Serialiser().Duration(t.Total)
		if t.Tag.Value() == "" {
			table.CellL("#" + t.Tag.Name())
			table.CellL(totalString)
		} else {
			if opt.Values {
				table.CellL(" " + t.Tag.Value())
				table.Skip(1)
				table.CellL(totalString)
			}
		}
		if t.Tag.Value() == "" && opt.Values {
			table.Skip(1)
		}
	}
	table.Collect(ctx.Print)
	opt.WarnArgs.PrintWarnings(ctx, records)
	return nil
}
