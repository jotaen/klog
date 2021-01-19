package cli

import (
	"fmt"
	"klog/app"
	"klog/parser"
	"klog/service"
)

type Print struct {
	FilterArgs
	FileArgs
	Sort bool `short:"s" name:"sort" help:"Sort output by date (from oldest to latest)"`
}

func (args *Print) Run(ctx *app.Context) error {
	rs, err := ctx.RetrieveRecords(args.File)
	if err != nil {
		return prettifyError(err)
	}
	rs, _ = service.FindFilter(rs, args.FilterArgs.toFilter())
	if args.Sort {
		service.Sort(rs)
	}
	fmt.Println("\n" + parser.SerialiseRecords(rs, styler))
	return nil
}
