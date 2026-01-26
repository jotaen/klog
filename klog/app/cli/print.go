package cli

import (
	"strings"

	"github.com/jotaen/klog/klog"
	"github.com/jotaen/klog/klog/app"
	"github.com/jotaen/klog/klog/app/cli/args"
	"github.com/jotaen/klog/klog/parser"
	"github.com/jotaen/klog/klog/service"
	tf "github.com/jotaen/klog/lib/terminalformat"
)

type Print struct {
	WithTotals bool `name:"with-totals" help:"Amend output with evaluated total times."`
	args.FilterArgs
	args.SortArgs
	args.WarnArgs
	args.NoStyleArgs
	args.InputFilesArgs
}

func (opt *Print) Help() string {
	return `
Outputs data on the terminal, by default with syntax-highlighting turned on.
Note that the output doesn’t resemble the file byte by byte, but the command may apply some minor clean-ups of the formatting.

If run with filter flags, it only outputs those entries that match the filter clauses.
You can optionally also sort the records, or print out the total times for each record and entry.
`
}

func (opt *Print) Run(ctx app.Context) app.Error {
	opt.NoStyleArgs.Apply(&ctx)
	styler, serialser := ctx.Serialise()
	records, err := ctx.ReadInputs(opt.File...)
	if err != nil {
		return err
	}
	now := ctx.Now()
	records, fErr := opt.ApplyFilter(now, records)
	if fErr != nil {
		return fErr
	}
	if len(records) == 0 {
		return nil
	}
	records = opt.ApplySort(records)
	serialisedRecords := parser.SerialiseRecords(serialser, records...)
	output := func() string {
		if opt.WithTotals {
			return printWithDurations(styler, serialisedRecords)
		}
		return "\n" + serialisedRecords.ToString()
	}()
	ctx.Print(output + "\n")

	opt.WarnArgs.PrintWarnings(ctx, records, nil)
	return nil
}

func printWithDurations(styler tf.Styler, ls parser.Lines) string {
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
				value += styler.Props(tf.StyleProps{Color: tf.TEXT_SUBDUED}).Format(p.d.ToString())
			} else {
				value += styler.Props(tf.StyleProps{IsUnderlined: true}).Format(p.d.ToString())
			}
			return strings.Repeat(" ", maxColumnLength-length+1) + value
		}()
		result += "  |  "
		result += l.Text
		result += "\n"
	}
	return result
}
