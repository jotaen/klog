package cli

import (
	"klog/app"
	"klog/parser"
	"klog/service"
)

type Print struct {
	FilterArgs
	SortArgs
	InputFilesArgs
}

func (args *Print) Run(ctx app.Context) error {
	rs, err := ctx.RetrieveRecords(args.File...)
	if err != nil {
		return err
	}
	if len(rs) == 0 {
		return nil
	}
	opts := args.FilterArgs.toFilter()
	opts.Sort = args.Sort
	rs = service.Query(rs, opts)
	ctx.Print("\n" + parser.SerialiseRecords(&styler, rs...) + "\n")
	return nil
}
