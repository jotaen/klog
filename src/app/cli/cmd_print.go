package cli

import (
	"klog/app"
	"klog/parser"
	"klog/service"
)

type Print struct {
	FilterArgs
	InputFilesArgs
	Sort bool `short:"s" name:"sort" help:"Sort output by date (from oldest to latest)"`
}

func (args *Print) Run(ctx app.Context) error {
	rs, err := ctx.RetrieveRecords(args.File...)
	if err != nil {
		return err
	}
	if len(rs) == 0 {
		return nil
	}
	rs = service.FindFilter(rs, args.FilterArgs.toFilter())
	if args.Sort {
		rs = service.Sort(rs, true)
	}
	ctx.Print("\n" + parser.SerialiseRecords(&styler, rs...) + "\n")
	return nil
}
