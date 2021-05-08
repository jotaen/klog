package cli

import (
	"klog/app"
	"klog/app/cli/lib"
	"klog/parser/json"
	"klog/parser/parsing"
)

type Json struct {
	lib.FilterArgs
	lib.SortArgs
	lib.InputFilesArgs
	Pretty bool `name:"pretty" help:"Pretty-print output"`
}

func (opt *Json) Help() string {
	return `Run with the --pretty flag to explore how the output structure looks.`
}

func (opt *Json) Run(ctx app.Context) error {
	records, err := ctx.ReadInputs(opt.File...)
	if err != nil {
		parserErrs, isParserErr := err.(parsing.Errors)
		if isParserErr {
			ctx.Print(json.ToJson(nil, parserErrs, opt.Pretty) + "\n")
			return nil
		}
		return err
	}
	records = opt.ApplyFilter(ctx.Now(), records)
	records = opt.ApplySort(records)
	ctx.Print(json.ToJson(records, nil, opt.Pretty) + "\n")
	return nil
}
