package cli

import (
	"klog/app"
	"klog/parser"
)

type Print struct {
	FilterArgs
	SortArgs
	WarnArgs
	InputFilesArgs
}

func (opt *Print) Run(ctx app.Context) error {
	records, err := ctx.RetrieveRecords(opt.File...)
	if err != nil {
		return err
	}
	if len(records) == 0 {
		return nil
	}
	now := ctx.Now()
	records = opt.filter(now, records)
	records = opt.sort(records)
	ctx.Print("\n" + parser.SerialiseRecords(&styler, records...) + "\n")

	ctx.Print(opt.WarnArgs.ToString(now, records))
	return nil
}
