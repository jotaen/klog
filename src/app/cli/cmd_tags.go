package cli

import (
	. "klog"
	"klog/app"
	"klog/service"
	"sort"
	"strings"
)

type Tags struct {
	FilterArgs
	WarnArgs
	InputFilesArgs
}

func (opt *Tags) Run(ctx app.Context) error {
	records, err := ctx.RetrieveRecords(opt.File...)
	if err != nil {
		return err
	}
	records = opt.filter(ctx.Now(), records)
	entriesByTag, _ := service.EntryTagLookup(records...)
	tagsOrdered, maxLength := sortTags(entriesByTag)
	for _, t := range tagsOrdered {
		es := entriesByTag[t]
		ctx.Print(t.ToString())
		ctx.Print(strings.Repeat(" ", maxLength-len(t)) + " ")
		ctx.Print(styler.Duration(service.TotalEntries(es...), false))
		ctx.Print("\n")
	}

	ctx.Print(opt.WarnArgs.ToString(ctx.Now(), records))
	return nil
}

func sortTags(ts map[Tag][]Entry) ([]Tag, int) {
	var result []Tag
	maxLength := 0
	for t := range ts {
		result = append(result, t)
		if len(t) > maxLength {
			maxLength = len(t)
		}
	}
	sort.Slice(result, func(i int, j int) bool {
		return result[i] < result[j]
	})
	return result, maxLength
}
