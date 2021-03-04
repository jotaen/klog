package cli

import (
	"klog/app"
	"klog/parser/json"
	"klog/parser/parsing"
)

type Json struct {
	FilterArgs
	SortArgs
	InputFilesArgs
}

func (opt *Json) Run(ctx app.Context) error {
	records, err := ctx.RetrieveRecords()
	if err != nil {
		parserErrs, isParserErr := err.(parsing.Errors)
		if isParserErr {
			ctx.Print(json.ToJson(nil, parserErrs) + "\n")
			return nil
		}
		return err
	}
	records = opt.filter(ctx.Now(), records)
	records = opt.sort(records)
	ctx.Print(json.ToJson(records, nil) + "\n")
	return nil
}
