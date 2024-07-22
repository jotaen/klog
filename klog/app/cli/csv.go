package cli

import (
	"fmt"
	"strings"

	"github.com/jotaen/klog/klog/app"
	"github.com/jotaen/klog/klog/app/cli/util"
)

type Csv struct {
	util.NowArgs
	util.FilterArgs
	util.SortArgs
	util.InputFilesArgs
}

func (opt *Csv) Help() string {
	return `
This commands outputs the records into a simple csv format with the following collumns:
| Date | Duration | Tags | Description |
Example:
date       ,duration, tag           , description          
2020-01-01 , 60     , #science      , Worked on the project
Please note: Entries with >1 tag will be repeated for each tag.
duration will always be in minutes.
If there are errors in the parsing, they will be printed on separate lines. and
the operation will return a non-zero exit code.
`
}

func (opt *Csv) Run(ctx app.Context) app.Error {
	records, err := ctx.ReadInputs(opt.File...)
	if err != nil {
		return err
	}
	fmt.Printf("date,duration,tag,description\n")
	for _, record := range records {
		entries := record.Entries()
		for _, entry := range entries {
			duration := entry.Duration().InMinutes()
			summary := entry.Summary()
			date := record.Date().ToString()

			text := strings.Join(summary.Lines(), " ")
			sanitizedText := []string{}
			words := strings.Fields(text)
			for _, word := range words {
				if !strings.HasPrefix(word, "#") {
					sanitizedText = append(sanitizedText, word)
				}
			}
			recordText := strings.Join(sanitizedText, " ")

			tags := summary.Tags().ToStrings()
			for _, tag := range tags {
				fmt.Printf("%s,%d,\"%s\",\"%s\"\n", date, duration, tag, recordText)
			}
		}
	}
	return nil
}
