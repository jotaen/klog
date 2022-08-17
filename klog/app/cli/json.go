package cli

import (
	"github.com/jotaen/klog/klog/app"
	"github.com/jotaen/klog/klog/app/cli/lib"
	"github.com/jotaen/klog/klog/parser/json"
)

type Json struct {
	lib.FilterArgs
	lib.SortArgs
	lib.InputFilesArgs
	Pretty bool `name:"pretty" help:"Pretty-print output"`
}

func (opt *Json) Help() string {
	return `The output structure contains two properties at the top level: "records" and "errors".

If the file is valid, "records" is an array containing a JSON object for each record; "errors" is null.

If the file has syntax errors, "records" is null and "errors" contains an array of error objects.

The structure of the objects is always uniform, so you can explore it by running the command with the --pretty flag.
`
}

func (opt *Json) Run(ctx app.Context) error {
	records, err := ctx.ReadInputs(opt.File...)
	if err != nil {
		parserErrs, isParserErr := err.(app.ParserErrors)
		if isParserErr {
			ctx.Print(json.ToJson(nil, parserErrs.All(), opt.Pretty) + "\n")
			return nil
		}
		return err
	}
	records = opt.ApplyFilter(ctx.Now(), records)
	records = opt.ApplySort(records)
	ctx.Print(json.ToJson(records, nil, opt.Pretty) + "\n")
	return nil
}
