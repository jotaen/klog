package cli

import (
	"fmt"
	. "klog"
	"klog/app"
	"klog/service"
)

type Total struct {
	FilterArgs
	DiffArg
	WarnArgs
	NowArgs
	InputFilesArgs
}

func (opt *Total) Run(ctx app.Context) error {
	records, err := ctx.RetrieveRecords(opt.File...)
	if err != nil {
		return err
	}
	records = opt.filter(ctx.Now(), records)
	total := opt.NowArgs.total(ctx.Now(), records...)
	ctx.Print(fmt.Sprintf("Total: %s\n", styler.Duration(total, false)))
	if opt.Diff {
		should := service.ShouldTotalSum(records...)
		diff := NewDuration(0, 0).Minus(should).Plus(total)
		ctx.Print(fmt.Sprintf("Should: %s\n", styler.ShouldTotal(should)))
		ctx.Print(fmt.Sprintf("Diff: %s\n", styler.Duration(diff, true)))
	}
	ctx.Print(fmt.Sprintf("(In %d record%s)\n", len(records), func() string {
		if len(records) == 1 {
			return ""
		}
		return "s"
	}()))

	ctx.Print(opt.WarnArgs.ToString(ctx.Now(), records))
	return nil
}
