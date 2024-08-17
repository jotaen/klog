package cli

import (
	"fmt"
	"github.com/jotaen/klog/klog/app"
	tf "github.com/jotaen/klog/klog/app/cli/terminalformat"
	"github.com/jotaen/klog/klog/app/cli/util"
	"github.com/jotaen/klog/klog/service"
)

type Tags struct {
	Values bool `name:"values" short:"v" help:"Display breakdown of tag values (if the data contains any; e.g.: '#tag=value')."`
	Count  bool `name:"count" short:"c" help:"Display the number of matching entries per tag."`
	util.FilterArgs
	util.NowArgs
	util.DecimalArgs
	util.WarnArgs
	util.NoStyleArgs
	util.InputFilesArgs
}

func (opt *Tags) Help() string {
	return ` 
If a tag appears in the overall record summary, then all of the record’s entries match.
If a tag appears in an entry summary, only that particular entry matches.
If tags are specified redundantly in the data, the respective time is still counted uniquely.

If you use tags with values (e.g., '#tag=value'), then these also match against the base tag (e.g., '#tag').
You can use the '--values' flag to display an additional breakdown by tag value.

Note that tag names are case-insensitive (e.g., '#tag' is the same as '#TAG'), whereas tag values are case-sensitive (so '#tag=value' is different from '#tag=VALUE').
`
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
	table := tf.NewTable(numberOfColumns, " ")
	for _, t := range totalByTag {
		totalString := serialiser.Duration(t.Total)
		countString := styler.Props(tf.StyleProps{Color: tf.TEXT_SUBDUED}).Format(fmt.Sprintf(" (%d)", t.Count))
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
			table.CellL(" " + styler.Props(tf.StyleProps{Color: tf.TEXT_SUBDUED}).Format(t.Tag.Value()))
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
