package cli

import (
	"klog/app"
	"klog/parser"
	"klog/service"
)

type Print struct {
	FilterArgs
	SortArgs
	WarnArgs
	InputFilesArgs
}

func (args *Print) Run(ctx app.Context) error {
	records, err := ctx.RetrieveRecords(args.File...)
	if err != nil {
		return err
	}
	if len(records) == 0 {
		return nil
	}
	opts := args.FilterArgs.toFilter()
	opts.Sort = args.Sort
	records = service.Query(records, opts)
	ctx.Print("\n" + parser.SerialiseRecords(&styler, records...) + "\n")

	args.WarnArgs.printWarnings(ctx, records)
	return nil
}
