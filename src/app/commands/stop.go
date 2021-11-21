package commands

import (
	. "github.com/jotaen/klog/src"
	"github.com/jotaen/klog/src/app"
	"github.com/jotaen/klog/src/app/lib"
	"github.com/jotaen/klog/src/parser/lineparsing"
	"github.com/jotaen/klog/src/parser/reconciler"
)

type Stop struct {
	lib.AtTimeArgs
	lib.AtDateArgs
	Summary string `name:"summary" short:"s" help:"Text to append to the entry summary"`
	lib.NoStyleArgs
	lib.OutputFileArgs
}

func (opt *Stop) Help() string {
	return `If the record contains an open-ended time range (e.g. 18:00-?) then this command
will replace the end placeholder with the current time (or the one specified via --time).`
}

func (opt *Stop) Run(ctx app.Context) error {
	opt.NoStyleArgs.Apply(&ctx)
	date := opt.AtDate(ctx.Now())
	time := opt.AtTime(ctx.Now())
	return ctx.ReconcileFile(opt.OutputFileArgs.File,
		func(records []Record, blocks []lineparsing.Block) (*reconciler.ReconcileResult, error) {
			entryReconciler := reconciler.NewEntryReconciler(records, blocks, func(r Record) bool {
				return r.Date().IsEqualTo(date)
			})
			if entryReconciler == nil {
				return nil, reconciler.NotEligibleError{}
			}
			return entryReconciler.CloseOpenRange(
				func(r Record) (Time, EntrySummary) { return time, NewEntrySummary(opt.Summary) },
			)
		},
		func(record []Record, blocks []lineparsing.Block) (*reconciler.ReconcileResult, error) {
			entryReconciler := reconciler.NewEntryReconciler(record, blocks, func(r Record) bool {
				return r.Date().IsEqualTo(date.PlusDays(-1))
			})
			if entryReconciler == nil {
				return nil, reconciler.NotEligibleError{}
			}
			adjustedTime := func() Time {
				if time.IsTomorrow() {
					return time
				}
				timeTomorrow, _ := time.Add(NewDuration(24, 0))
				return timeTomorrow
			}()
			return entryReconciler.CloseOpenRange(
				func(r Record) (Time, EntrySummary) { return adjustedTime, NewEntrySummary(opt.Summary) },
			)
		},
	)
}
