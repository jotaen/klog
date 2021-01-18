package cli

import (
	"fmt"
	"klog/app"
	"klog/record"
	"klog/service"
)

type Total struct {
	FilterArgs
	FileArgs
	Diff bool `name:"diff" help:"Show diff between should and actual time"`
}

func (args *Total) Run(ctx *app.Context) error {
	rs, err := retrieveRecords(ctx, args.File)
	if err != nil {
		return err
	}
	rs, es := service.FindFilter(rs, args.FilterArgs.toFilter())
	total := service.TotalEntries(es)
	fmt.Printf("Total: %s\n", styler.PrintDuration(total))
	if args.Diff {
		should := service.ShouldTotalAll(rs)
		diff := record.NewDuration(0, 0).Subtract(should).Add(total)
		fmt.Printf("Should: %s\n", styler.PrintDuration(should))
		fmt.Printf("Diff: %s\n", styler.PrintDuration(diff))
	}
	fmt.Printf("(In %d records)\n", len(rs))
	return nil
}
