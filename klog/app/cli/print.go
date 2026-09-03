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
	args.NowArgs
	args.WarnArgs
	args.NoStyleArgs
	args.InputFilesArgs
}

func (opt *Print) Help() string {
	return `
Outputs data on the terminal, by default with syntax-highlighting turned on.
Note that the output doesn’t resemble the file verbatim, but it may apply some minor formatting.

You can optionally also sort the records, or print out the total times for each record and entry.

Open-ended time ranges (e.g., '8:00 - ?') are printed as they are, and they count as '0m' in the totals.
With the '--now' flag, they are printed as if they were closed “right now”, and their elapsed time is factored into the totals.
In the totals column, such an entry duration is put in parentheses, e.g. '(2h17m)'.
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
	closedByNow := map[klog.Record]int{}
	if opt.Now {
		closedByNow = openEntryIndices(records)
	}
	nErr := opt.ApplyNow(now, records...)
	if nErr != nil {
		return nErr
	}
	serialisedRecords := parser.SerialiseRecords(serialser, records...)
	output := func() string {
		if opt.WithTotals {
			return printWithDurations(styler, serialisedRecords, closedByNow)
		}
		return "\n" + serialisedRecords.ToString()
	}()
	ctx.Print(output + "\n")

	opt.WarnArgs.PrintWarnings(ctx, records, []service.UsageWarning{opt.NowArgs.GetWarning()})
	return nil
}

// openEntryIndices maps every record that has an open-ended time range to
// the index of that entry. Records without an open range are not included.
func openEntryIndices(rs []klog.Record) map[klog.Record]int {
	result := map[klog.Record]int{}
	for _, r := range rs {
		for entryI, e := range r.Entries() {
			isOpen := klog.Unbox[bool](&e,
				func(klog.Range) bool { return false },
				func(klog.Duration) bool { return false },
				func(klog.OpenRange) bool { return true },
			)
			if isOpen {
				result[r] = entryI
				break
			}
		}
	}
	return result
}

// printWithDurations prefixes each line with its total duration. `closedByNow`
// denotes the entries whose open range was closed via `--now`; their duration
// is put in parentheses, as it depends on the current time.
func printWithDurations(styler tf.Styler, ls parser.Lines, closedByNow map[klog.Record]int) string {
	type Prefix struct {
		d           klog.Duration
		isSub       bool
		closedByNow bool
	}
	text := func(p *Prefix) string {
		if p.closedByNow {
			return "(" + p.d.ToString() + ")"
		}
		return p.d.ToString()
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
				return &Prefix{service.Total(l.Record), false, false}
			}
			if l.EntryI != -1 && l.EntryI != previousEntry {
				previousEntry = l.EntryI
				openI, hasOpen := closedByNow[l.Record]
				return &Prefix{l.Record.Entries()[l.EntryI].Duration(), true, hasOpen && openI == l.EntryI}
			} else {
				return nil
			}
		}()
		prefixes = append(prefixes, prefix)
		if prefix != nil && len(text(prefix)) > maxColumnLength {
			maxColumnLength = len(text(prefix))
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
			length := len(text(p))
			value := ""
			if p.isSub {
				value += styler.Props(tf.StyleProps{Color: tf.TEXT_SUBDUED}).Format(text(p))
			} else {
				value += styler.Props(tf.StyleProps{IsUnderlined: true}).Format(text(p))
			}
			return strings.Repeat(" ", maxColumnLength-length+1) + value
		}()
		result += "  |  "
		result += l.Text
		result += "\n"
	}
	return result
}
