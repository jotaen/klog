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
	lib.InputFilesArgs
}

func (opt *Total) Run(ctx app.Context) error {
	records, err := ctx.ReadInputs(opt.File...)
	if err != nil {
		return err
	}
	now := ctx.Now()
	records = opt.ApplyFilter(now, records)
	total := opt.NowArgs.Total(now, records...)
	ctx.Print(fmt.Sprintf("Total: %s\n", lib.Styler.Duration(total, false)))
	if opt.Diff {
		should := service.ShouldTotalSum(records...)
		diff := service.Diff(should, total)
		ctx.Print(fmt.Sprintf("Should: %s\n", lib.Styler.ShouldTotal(should)))
		ctx.Print(fmt.Sprintf("Diff: %s\n", lib.Styler.Duration(diff, true)))
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
