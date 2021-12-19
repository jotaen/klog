package commands

import (
	. "github.com/jotaen/klog/src"
	"github.com/jotaen/klog/src/app"
	"github.com/jotaen/klog/src/app/lib"
	"github.com/jotaen/klog/src/parser"
	"github.com/jotaen/klog/src/parser/reconciling"
)

type Stop struct {
	lib.AtDateAndTimeArgs
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
	now := ctx.Now()
	date, isAutoDate := opt.AtDate(now)
	time, isAutoTime, err := opt.AtTime(now)
	if err != nil {
		return err
	}
	return ctx.ReconcileFile(
		opt.OutputFileArgs.File,

		[]reconciling.Creator{
			func(parsedRecords []parser.ParsedRecord) *reconciling.Reconciler {
				return reconciling.NewReconcilerAtRecord(parsedRecords, date)
			},
			func(parsedRecords []parser.ParsedRecord) *reconciling.Reconciler {
				if isAutoDate && isAutoTime {
					// Only fall back to yesterday if no explicit date has been given.
					// Otherwise, it wouldn’t make sense to decrement the day.
					time, _ = time.Add(NewDuration(24, 0))
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
