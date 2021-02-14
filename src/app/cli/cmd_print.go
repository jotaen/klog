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
	records = opt.filter(ctx.Now(), records)
	records = opt.sort(records)
	ctx.Print("\n" + parser.SerialiseRecords(&styler, records...) + "\n")

	ctx.Print(opt.WarnArgs.ToString(ctx.Now(), records))
	return nil
}
