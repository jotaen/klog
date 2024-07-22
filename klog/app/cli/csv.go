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
This command outputs the records into a simple csv format with the following columns:

Example:
date, duration, tag, description
2020-01-01 , 60     , #science      , Worked on the project
Please note: Entries with >1 tag will be repeated for each tag.
duration will always be in minutes.
If there are errors in parsing, they will be printed on separate lines, and
the operation will return a non-zero exit code.
`
}

const csvHeader = "date,duration,tag,description\n"

func sanitizeText(text string) string {
	sanitizedText := []string{}
	words := strings.Fields(text)
	for _, word := range words {
		if !strings.HasPrefix(word, "#") {
			sanitizedText = append(sanitizedText, word)
		}
	}
	return strings.Join(sanitizedText, " ")
}

func (opt *Csv) Run(ctx app.Context) app.Error {
	records, err := ctx.ReadInputs(opt.File...)
	if err != nil {
		return err
	}
	fmt.Print(csvHeader)
	for _, record := range records {
		entries := record.Entries()
		if len(entries) == 0 {
			fmt.Printf("%s,0,,\n", record.Date().ToString())
			continue
		}
		for _, entry := range entries {
			date := record.Date().ToString()
			duration := entry.Duration().InMinutes()
			summary := entry.Summary()
			recordText := sanitizeText(strings.Join(summary.Lines(), " "))

			tags := summary.Tags().ToStrings()
			if len(tags) == 0 {
				fmt.Printf("%s,%d,,\"%s\"\n", date, duration, recordText)
			} else {
				for _, tag := range tags {
					fmt.Printf("%s,%d,\"%s\",\"%s\"\n", date, duration, tag, recordText)
				}
			}
		}
	}
	return nil
}
