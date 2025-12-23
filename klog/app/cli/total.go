package cli

import (
	"fmt"

	"github.com/jotaen/klog/klog/app"
	"github.com/jotaen/klog/klog/app/cli/util"
	"github.com/jotaen/klog/klog/service"
)

type Total struct {
	util.FilterArgs
	util.DiffArgs
	util.NowArgs
	util.DecimalArgs
	util.WarnArgs
	util.NoStyleArgs
	util.InputFilesArgs
}

func (opt *Total) Help() string {
	return `
By default, the total time consists of all durations and time ranges, but it doesn’t include open-ended time ranges (e.g., '8:00 - ?').
If you want to factor them in anyway, you can use the '--now' option, which treats all open-ended time ranges as if they were closed “right now”.

If the records contain should-total values, you can also compute the difference between should-total and actual total by using the '--diff' flag.
`
}

func (opt *Total) Run(ctx app.Context) app.Error {
	opt.DecimalArgs.Apply(&ctx)
	opt.NoStyleArgs.Apply(&ctx)
	_, serialiser := ctx.Serialise()
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
	total := service.Total(records...)
	ctx.Print(fmt.Sprintf("Total: %s\n", serialiser.Duration(total)))
	if opt.Diff {
		should := service.ShouldTotalSum(records...)
		diff := service.Diff(should, total)
		ctx.Print(fmt.Sprintf("Should: %s\n", serialiser.ShouldTotal(should)))
		ctx.Print(fmt.Sprintf("Diff: %s\n", serialiser.SignedDuration(diff)))
	}
	ctx.Print(fmt.Sprintf("(In %d record%s)\n", len(records), func() string {
		if len(records) == 1 {
			return ""
		}
		return "s"
	}()))

	opt.WarnArgs.PrintWarnings(ctx, records, opt.GetNowWarnings())
	return nil
}
