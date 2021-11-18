package cli

import (
	"github.com/jotaen/klog/lib/jotaen/terminalformat"
	. "github.com/jotaen/klog/src"
	"github.com/jotaen/klog/src/app"
	"github.com/jotaen/klog/src/app/cli/lib"
	"github.com/jotaen/klog/src/service"
	"sort"
)

type Tags struct {
	lib.FilterArgs
	lib.WarnArgs
	lib.NoStyleArgs
	lib.InputFilesArgs
}

func (opt *Tags) Run(ctx app.Context) error {
	opt.NoStyleArgs.Apply(&ctx)
	records, err := ctx.ReadInputs(opt.File...)
	if err != nil {
		return err
	}
	now := ctx.Now()
	records = opt.ApplyFilter(now, records)
	entriesByTag := service.EntryTagLookup(records...)
	tagsOrdered := sortTags(entriesByTag)
	if len(tagsOrdered) == 0 {
		return nil
	}
	table := terminalformat.NewTable(2, " ")
	for _, t := range tagsOrdered {
		es := entriesByTag[t]
		table.
			CellL(t.ToString()).
			CellL(ctx.Serialiser().Duration(service.TotalEntries(es...)))
	}
	table.Collect(ctx.Print)
	ctx.Print(opt.WarnArgs.ToString(now, records))
	return nil
}

func sortTags(ts map[Tag][]Entry) []Tag {
	var result []Tag
	for t := range ts {
		result = append(result, t)
	}
	sort.Slice(result, func(i int, j int) bool {
		return result[i] < result[j]
	})
	return result
}
