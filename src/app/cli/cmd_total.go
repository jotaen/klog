package cli

import (
	"errors"
	"fmt"
	"klog/app"
	"klog/service"
)

type Total struct {
	FilterArgs
	FileArgs
}

func (args *Total) Run(ctx *app.Context) error {
	rs, err := ctx.Read(args.File)
	if err != nil {
		return errors.New("EXECUTION_FAILED")
	}
	rs, es := service.FindFilter(rs, args.FilterArgs.ToFilter())
	total := service.TotalEntries(es)
	fmt.Printf("Total: %s\n", total.ToString())
	fmt.Printf("(In %d records)\n", len(rs))
	return nil
}
