package cli

import (
	"github.com/jotaen/klog/klog"
	"github.com/jotaen/klog/klog/app"
	"github.com/jotaen/klog/klog/app/cli/lib"
	"github.com/jotaen/klog/klog/app/cli/lib/terminalformat"
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
	return `The output is syntax-highlighted. Note that the formatting is sanitised/normalised, especially in regards to whitespace.`
}

func (opt *Print) Run(ctx app.Context) app.Error {
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

	opt.WarnArgs.PrintWarnings(ctx, records, nil)
	return nil
}

func printWithDurations(serialiser parser.Serialiser, ls parser.Lines) string {
	type Prefix struct {
		d     klog.Duration
		isSub bool
	}
	var prefixes []*Prefix
	maxColumnLength := 0
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
				return &Prefix{service.Total(l.Record), false}
			}
			if l.EntryI != -1 && l.EntryI != previousEntry {
				previousEntry = l.EntryI
				return &Prefix{l.Record.Entries()[l.EntryI].Duration(), true}
			} else {
				return nil
			}
		}()
		prefixes = append(prefixes, prefix)
		if prefix != nil && len(prefix.d.ToString()) > maxColumnLength {
			maxColumnLength = len(prefix.d.ToString())
		}
	}

	result := "\n"
	for i, l := range ls {
		p := prefixes[i]
		if l.Record == nil {
			result += "\n"
			continue
		}
		result += func() string {
			if p == nil {
				return strings.Repeat(" ", maxColumnLength+1)
			}
			length := len(p.d.ToString())
			value := ""
			if p.isSub {
				value += serialiser.Format(terminalformat.Style{Color: "247"}, p.d.ToString())
			} else {
				value += serialiser.Format(terminalformat.Style{IsUnderlined: true}, p.d.ToString())
			}
			return strings.Repeat(" ", maxColumnLength-length+1) + value
		}()
		result += "  |  "
		result += l.Text
		result += "\n"
	}
	return result
}
