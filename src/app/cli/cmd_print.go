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
}

func (args *Print) Run(ctx *app.Context) error {
	rs, err := retrieveRecords(ctx, args.File)
	if err != nil {
		return err
	}
	rs, _ = service.FindFilter(rs, args.FilterArgs.toFilter())
	h := cliPrinter{}
	fmt.Println(parser.SerialiseRecords(rs, h))
	return nil
}
