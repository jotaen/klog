package cli

import (
	"fmt"

	"github.com/jotaen/klog/klog/app"
	"github.com/jotaen/klog/klog/app/cli/args"
	"github.com/jotaen/klog/klog/service"
	tf "github.com/jotaen/klog/lib/terminalformat"
)

type Tags struct {
	Values       bool `name:"values" short:"v" help:"Display breakdown of tag values (if the data contains any; e.g.: '#tag=value')."`
	Count        bool `name:"count" short:"c" help:"Display the number of matching entries per tag."`
	WithUntagged bool `name:"with-untagged" short:"u" help:"Display remainder of any untagged entries"`
	args.FilterArgs
	args.NowArgs
	args.DecimalArgs
	args.WarnArgs
	args.NoStyleArgs
	args.InputFilesArgs
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
	records, fErr := opt.ApplyFilter(now, records)
	if fErr != nil {
		return fErr
	}
	nErr := opt.ApplyNow(now, records...)
	if nErr != nil {
		return nErr
	}
	tagStats, untagged := service.AggregateTotalsByTags(records...)
	numberOfColumns := 2
	if opt.Values {
		numberOfColumns++
	}
	if opt.Count {
		numberOfColumns++
	}
	countString := func(c int) string {
		return styler.Props(tf.StyleProps{Color: tf.TEXT_SUBDUED}).Format(fmt.Sprintf(" (%d)", c))
	}
	table := tf.NewTable(numberOfColumns, " ")
	for _, t := range tagStats {
		totalString := serialiser.Duration(t.Total)
		if t.Tag.Value() == "" {
			table.CellL("#" + t.Tag.Name())
			table.CellL(totalString)
			if opt.Values {
				table.Skip(1)
			}
			if opt.Count {
				table.CellL(countString(t.Count))
			}
		} else if opt.Values {
			table.CellL(" " + styler.Props(tf.StyleProps{Color: tf.TEXT_SUBDUED}).Format(t.Tag.Value()))
			table.Skip(1)
			table.CellL(totalString)
			if opt.Count {
				table.CellL(countString(t.Count))
			}
		}
	}
	if opt.WithUntagged {
		table.CellL("(untagged)")
		table.CellL(serialiser.Duration(untagged.Total))
		if opt.Values {
			table.Skip(1)
		}
		if opt.Count {
			table.CellL(countString(untagged.Count))
		}
	}
	table.Collect(ctx.Print)
	opt.WarnArgs.PrintWarnings(ctx, records, []service.UsageWarning{opt.NowArgs.GetWarning()})
	return nil
}
