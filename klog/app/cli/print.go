package cli

import (
	"github.com/jotaen/klog/klog"
	"github.com/jotaen/klog/klog/app"
	"github.com/jotaen/klog/klog/app/cli/lib"
	"github.com/jotaen/klog/klog/parser"
	"github.com/jotaen/klog/klog/service"
	"strings"
)

type Print struct {
	WithTotals bool `name:"with-totals" help:"Amend output with evaluated total times"`
	lib.FilterArgs
	lib.SortArgs
	lib.WarnArgs
	lib.NoStyleArgs
	lib.InputFilesArgs
}

func (opt *Print) Help() string {
	return `The output is syntax-highlighted and the formatting is slightly sanitised.`
}

func (opt *Print) Run(ctx app.Context) error {
	opt.NoStyleArgs.Apply(&ctx)
	records, err := ctx.ReadInputs(opt.File...)
	if err != nil {
		return err
	}
	now := ctx.Now()
	records = opt.ApplyFilter(now, records)
	if len(records) == 0 {
		return nil
	}
	records = opt.ApplySort(records)
	serialisedRecords := parser.SerialiseRecords(ctx.Serialiser(), records...)
	output := func() string {
		if opt.WithTotals {
			return printWithDurations(ctx.Serialiser(), serialisedRecords)
		}
		return "\n" + serialisedRecords.ToString()
	}()
	ctx.Print(output + "\n")

	opt.WarnArgs.PrintWarnings(ctx, records)
	return nil
}

func printWithDurations(serialiser parser.Serialiser, ls parser.Lines) string {
	type Prefix struct {
		d      klog.Duration
		column int
	}
	var prefixes []*Prefix
	maxColumnLengths := []int{0, 0}
	var previousRecord klog.Record
	previousEntry := -1
	for _, l := range ls {
		prefix := func() *Prefix {
			if l.Record == nil {
				previousRecord = nil
				previousEntry = -1
				return nil
			}
			if previousRecord == nil {
				previousRecord = l.Record
				return &Prefix{service.Total(l.Record), 0}
			}
			if l.EntryI != -1 && l.EntryI != previousEntry {
				previousEntry = l.EntryI
				return &Prefix{l.Record.Entries()[l.EntryI].Duration(), 1}
			} else {
				return nil
			}
		}()
		prefixes = append(prefixes, prefix)
		if prefix != nil && len(prefix.d.ToString()) > maxColumnLengths[prefix.column] {
			maxColumnLengths[prefix.column] = len(prefix.d.ToString())
		}
	}
	RECORD_SEPARATOR := strings.Repeat("-", maxColumnLengths[0]) + "-+-" + strings.Repeat("-", maxColumnLengths[1])
	result := RECORD_SEPARATOR + "-+ " + "\n"
	for i, l := range ls {
		prefixText := ""
		p := prefixes[i]
		if l.Record == nil {
			prefixText = RECORD_SEPARATOR
			prefixText += "-+ "
		} else {
			column := []string{strings.Repeat(" ", maxColumnLengths[0]), strings.Repeat(" ", maxColumnLengths[1])}
			if p != nil {
				column[p.column] = strings.Repeat(" ", maxColumnLengths[0]-len(p.d.ToString()))
				column[p.column] += serialiser.Duration(p.d)
			}
			prefixText += column[0]
			prefixText += " | "
			prefixText += column[1]
			prefixText += " | "
		}
		result += prefixText
		result += l.Text
		result += "\n"
	}
	result += RECORD_SEPARATOR + "-+ "
	return result
}
