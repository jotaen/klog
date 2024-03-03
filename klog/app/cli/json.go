package cli

import (
	"github.com/jotaen/klog/klog/app"
	"github.com/jotaen/klog/klog/app/cli/util"
	"github.com/jotaen/klog/klog/parser/json"
)

type Json struct {
	Pretty bool `name:"pretty" help:"Pretty-print output."`
	util.NowArgs
	util.FilterArgs
	util.SortArgs
	util.InputFilesArgs
}

func (opt *Json) Help() string {
	return `
The output structure is a JSON object which contains two properties at the top level: 'records' and 'errors'.
If the file is valid, 'records' is an array containing a JSON object for each record, and 'errors' is 'null'.
If the file has syntax errors, 'records' is 'null', and 'errors' contains an array of error objects.

The structure of the 'record' and 'error' objects is always uniform and should be self-explanatory.
You can best explore it by running the command with the --pretty flag.
`
}

func (opt *Json) Run(ctx app.Context) app.Error {
	records, err := ctx.ReadInputs(opt.File...)
	if err != nil {
		parserErrs, isParserErr := err.(app.ParserErrors)
		if isParserErr {
			ctx.Print(json.ToJson(nil, parserErrs.All(), opt.Pretty) + "\n")
			return nil
		}
		return err
	}
	now := ctx.Now()
	nErr := opt.ApplyNow(now, records...)
	if nErr != nil {
		return nErr
	}
	records = opt.ApplyFilter(now, records)
	records = opt.ApplySort(records)
	ctx.Print(json.ToJson(records, nil, opt.Pretty) + "\n")
	return nil
}
