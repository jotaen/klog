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

Open-ended time ranges (e.g., '8:00 - ?') are printed as they are, and they count as '0m' in the totals when using the '--with-totals' flag.
With the '--now' flag, they are printed as if they were closed “right now”, and their elapsed time is factored into the totals.
In the totals column, the provisional durations are marked with an asterisk, e.g. '2h17m*'.
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
	closedEntries, nErr := opt.ApplyNow(now, records...)
	if nErr != nil {
		return nErr
	}
	serialisedRecords := parser.SerialiseRecords(serialser, records...)
	output := func() string {
		if opt.WithTotals {
			return printWithDurations(styler, serialisedRecords, closedEntries)
		}
		return "\n" + serialisedRecords.ToString()
	}()
	ctx.Print(output + "\n")

	opt.WarnArgs.PrintWarnings(ctx, records, []service.UsageWarning{opt.NowArgs.GetWarning()})
	return nil
}

// printWithDurations prefixes each line with its total duration. `closedEntries`
// denotes the entries whose open range was closed via `--now`; the auto-closed durations
// are suffixed with an asterisk, to indicate that they are provisional.
func printWithDurations(styler tf.Styler, ls parser.Lines, closedEntries map[klog.Record]int) string {
	type Prefix struct {
		text  string
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
				return &Prefix{
					text: func() string {
						suffix := ""
						_, hasAutoClosedEntry := closedEntries[l.Record]
						if hasAutoClosedEntry {
							suffix = "*"
						}
						return service.Total(l.Record).ToString() + suffix
					}(),
					isSub: false,
				}
			}
			if l.EntryI != -1 && l.EntryI != previousEntry {
				previousEntry = l.EntryI
				return &Prefix{
					text: func() string {
						openI, hasAutoClosedEntry := closedEntries[l.Record]
						entryTotal := l.Record.Entries()[l.EntryI].Duration()
						suffix := ""
						if hasAutoClosedEntry && openI == l.EntryI {
							suffix = "*"
						}
						return entryTotal.ToString() + suffix
					}(),
					isSub: true,
				}
			} else {
				return nil
			}
		}()
		prefixes = append(prefixes, prefix)
		if prefix != nil && len(prefix.text) > maxColumnLength {
			maxColumnLength = len(prefix.text)
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
			length := len(p.text)
			value := ""
			if p.isSub {
				value += styler.Props(tf.StyleProps{Color: tf.TEXT_SUBDUED}).Format(p.text)
			} else {
				value += styler.Props(tf.StyleProps{IsUnderlined: true}).Format(p.text)
			}
			return strings.Repeat(" ", maxColumnLength-length+1) + value
		}()
		result += "  |  "
		result += l.Text
		result += "\n"
	}
	return result
}
