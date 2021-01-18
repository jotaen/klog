package cli

import (
	"fmt"
	"klog/app"
	"klog/service"
)

type Total struct {
	FilterArgs
	FileArgs
}

func (args *Total) Run(ctx *app.Context) error {
	rs, err := retrieveRecords(ctx, args.File)
	if err != nil {
		return err
	}
	rs, es := service.FindFilter(rs, args.FilterArgs.toFilter())
	total := service.TotalEntries(es)
	fmt.Printf("Total: %s\n", total.ToString())
	fmt.Printf("(In %d records)\n", len(rs))
	return nil
}
