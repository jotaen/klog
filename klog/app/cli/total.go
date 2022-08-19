package cli

import (
	"fmt"
	"github.com/jotaen/klog/klog/app"
	"github.com/jotaen/klog/klog/app/cli/lib"
	"github.com/jotaen/klog/klog/service"
)

type Total struct {
	lib.FilterArgs
	lib.DiffArgs
	lib.NowArgs
	lib.DecimalArgs
	lib.WarnArgs
	lib.NoStyleArgs
	lib.InputFilesArgs
}

func (opt *Total) Help() string {
	return `The total time is the sum of all records.

Note that the total time by default doesn’t include open-ended time ranges.
If you want to factor them in anyway, you can use the --now option,
which treats all open-ended time ranges as if they were closed right now.`
}

func (opt *Total) Run(ctx app.Context) error {
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
	total := service.Total(records...)
	ctx.Print(fmt.Sprintf("Total: %s\n", ctx.Serialiser().Duration(total)))
	if opt.Diff {
		should := service.ShouldTotalSum(records...)
		diff := service.Diff(should, total)
		ctx.Print(fmt.Sprintf("Should: %s\n", ctx.Serialiser().ShouldTotal(should)))
		ctx.Print(fmt.Sprintf("Diff: %s\n", ctx.Serialiser().SignedDuration(diff)))
	}
	ctx.Print(fmt.Sprintf("(In %d record%s)\n", len(records), func() string {
		if len(records) == 1 {
			return ""
		}
		return "s"
	}()))

	opt.WarnArgs.PrintWarnings(ctx, records)
	return nil
}
