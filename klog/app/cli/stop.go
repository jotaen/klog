package cli

import (
	"github.com/jotaen/klog/klog"
	"github.com/jotaen/klog/klog/app"
	"github.com/jotaen/klog/klog/app/cli/lib"
	"github.com/jotaen/klog/klog/parser"
	"github.com/jotaen/klog/klog/parser/reconciling"
)

type Stop struct {
	Summary klog.EntrySummary `name:"summary" short:"s" placeholder:"TEXT" help:"Text to append to the entry summary"`
	lib.AtDateAndTimeArgs
	lib.NoStyleArgs
	lib.OutputFileArgs
	lib.WarnArgs
}

func (opt *Stop) Help() string {
	return `If the record contains an open-ended time range (e.g. 18:00-?) then this command
will replace the end placeholder with the current time (or the one specified via --time).`
}

func (opt *Stop) Run(ctx app.Context) error {
	opt.NoStyleArgs.Apply(&ctx)
	now := ctx.Now()
	date, isAutoDate := opt.AtDate(now)
	time, isAutoTime, err := opt.AtTime(now)
	if err != nil {
		return err
	}
	return lib.Reconcile(ctx, lib.ReconcileOpts{OutputFileArgs: opt.OutputFileArgs, WarnArgs: opt.WarnArgs},
		[]reconciling.Creator{
			func(parsedRecords []parser.ParsedRecord) *reconciling.Reconciler {
				return reconciling.NewReconcilerAtRecord(parsedRecords, date)
			},
			func(parsedRecords []parser.ParsedRecord) *reconciling.Reconciler {
				if isAutoDate && isAutoTime {
					// Only fall back to yesterday if no explicit date has been given.
					// Otherwise, it wouldn’t make sense to decrement the day.
					time, _ = time.Plus(klog.NewDuration(24, 0))
					return reconciling.NewReconcilerAtRecord(parsedRecords, date.PlusDays(-1))
				}
				return nil
			},
		},

		func(reconciler *reconciling.Reconciler) (*reconciling.Result, error) {
			return reconciler.CloseOpenRange(time, opt.Summary)
		},
	)
}
