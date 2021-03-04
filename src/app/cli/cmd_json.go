package cli

import (
	"klog/app"
	"klog/parser/json"
	"klog/parser/parsing"
)

type Json struct {
	InputFilesArgs
}

func (opt *Json) Run(ctx app.Context) error {
	rs, err := ctx.RetrieveRecords()
	if err != nil {
		parserErrs, isParserErr := err.(parsing.Errors)
		if isParserErr {
			ctx.Print(json.ToJson(nil, parserErrs) + "\n")
			return nil
		}
		return err
	}
	ctx.Print(json.ToJson(rs, nil) + "\n")
	return nil
}
