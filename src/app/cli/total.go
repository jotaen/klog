package cli

import (
	"fmt"
	"klog/app"
	"klog/app/cli/lib"
	"klog/service"
)

type Total struct {
	lib.FilterArgs
	lib.DiffArgs
	lib.WarnArgs
	lib.NowArgs
	lib.NoStyleArgs
	lib.InputFilesArgs
}

func (opt *Total) Run(ctx app.Context) error {
	opt.NoStyleArgs.Apply(&ctx)
	records, err := ctx.ReadInputs(opt.File...)
	if err != nil {
		return err
	}
	now := ctx.Now()
	records = opt.ApplyFilter(now, records)
	total := opt.NowArgs.Total(now, records...)
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

	ctx.Print(opt.WarnArgs.ToString(now, records))
	return nil
}
