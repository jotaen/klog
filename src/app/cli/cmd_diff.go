package cli

import (
	"fmt"
	"klog/app"
	"klog/record"
	"klog/service"
)

type Diff struct {
	FilterArgs
	FilesArgs
}

func (args *Diff) Run(ctx *app.Context) error {
	rs, err := ctx.RetrieveRecords(args.File)
	if err != nil {
		return prettifyError(err)
	}
	rs, es := service.FindFilter(rs, args.FilterArgs.toFilter())
	total := service.TotalEntries(es)
	fmt.Printf("Total: %s\n", styler.PrintDuration(total))
	should := service.ShouldTotal(rs...)
	diff := record.NewDuration(0, 0).Minus(should).Plus(total)
	fmt.Printf("Should: %s\n", styler.PrintDuration(should))
	fmt.Printf("Diff: %s\n", styler.PrintDuration(diff))
	fmt.Printf("(In %d records)\n", len(rs))
	return nil
}
