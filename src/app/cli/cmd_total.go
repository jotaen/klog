package cli

import (
	"fmt"
	"klog/app"
	"klog/service"
)

type Total struct {
	FilterArgs
	FilesArgs
}

func (args *Total) Run(ctx *app.Context) error {
	rs, err := ctx.RetrieveRecords(args.File)
	if err != nil {
		return prettifyError(err)
	}
	rs, es := service.FindFilter(rs, args.FilterArgs.toFilter())
	total := service.TotalEntries(es)
	fmt.Printf("Total: %s\n", styler.PrintDuration(total))
	fmt.Printf("(In %d records)\n", len(rs))
	return nil
}
