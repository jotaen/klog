package cli

import (
	"fmt"
	"github.com/jotaen/klog/klog/app"
	terminalformat2 "github.com/jotaen/klog/klog/app/cli/terminalformat"
	"github.com/jotaen/klog/klog/app/cli/util"
	"github.com/jotaen/klog/klog/service"
)

type Tags struct {
	Values bool `name:"values" short:"v" help:"Display breakdown of tag values"`
	Count  bool `name:"count" short:"c" help:"Display the number of matching entries per tag"`
	util.FilterArgs
	util.NowArgs
	util.DecimalArgs
	util.WarnArgs
	util.NoStyleArgs
	util.InputFilesArgs
}

func (opt *Tags) Help() string {
	return `Aggregates the total times of entries by tags.

If a tag appears in the overall record summary, then all of the record’s entries match. If a tag appears in an entry summary, only that particular entry matches.

Every matching entry is counted individually.`
}

func (opt *Tags) Run(ctx app.Context) app.Error {
	opt.DecimalArgs.Apply(&ctx)
	opt.NoStyleArgs.Apply(&ctx)
	styler, serialiser := ctx.Serialise()
	records, err := ctx.ReadInputs(opt.File...)
	if err != nil {
		return err
	}
	now := ctx.Now()
	records = opt.ApplyFilter(now, records)
	nErr := opt.ApplyNow(now, records...)
	if nErr != nil {
		return nErr
	}
	totalByTag := service.AggregateTotalsByTags(records...)
	if len(totalByTag) == 0 {
		return nil
	}
	numberOfColumns := 2
	if opt.Values {
		numberOfColumns++
	}
	if opt.Count {
		numberOfColumns++
	}
	table := terminalformat2.NewTable(numberOfColumns, " ")
	for _, t := range totalByTag {
		totalString := serialiser.Duration(t.Total)
		countString := styler.Props(terminalformat2.StyleProps{Color: terminalformat2.SUBDUED}).Format(fmt.Sprintf(" (%d)", t.Count))
		if t.Tag.Value() == "" {
			table.CellL("#" + t.Tag.Name())
			table.CellL(totalString)
			if opt.Values {
				table.Skip(1)
			}
			if opt.Count {
				table.CellL(countString)
			}
		} else if opt.Values {
			table.CellL(" " + styler.Props(terminalformat2.StyleProps{Color: terminalformat2.SUBDUED}).Format(t.Tag.Value()))
			table.Skip(1)
			table.CellL(totalString)
			if opt.Count {
				table.CellL(countString)
			}
		}
	}
	table.Collect(ctx.Print)
	opt.WarnArgs.PrintWarnings(ctx, records, opt.GetNowWarnings())
	return nil
}
