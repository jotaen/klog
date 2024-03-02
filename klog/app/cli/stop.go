package cli

import (
	"github.com/jotaen/klog/klog"
	"github.com/jotaen/klog/klog/app"
	"github.com/jotaen/klog/klog/app/cli/util"
	"github.com/jotaen/klog/klog/parser/reconciling"
)

type Stop struct {
	Summary klog.EntrySummary `name:"summary" short:"s" placeholder:"TEXT" help:"Text to append to the entry summary"`
	util.AtDateAndTimeArgs
	util.NoStyleArgs
	util.OutputFileArgs
	util.WarnArgs
}

func (opt *Stop) Help() string {
	return `If the record contains an open-ended time range (e.g. 18:00-?) then this command
will replace the end placeholder with the current time (or the one specified via --time).

If the --time flag is not specified, it defaults to the current time as end time. In the latter case, the time can be rounded via --round.`
}

func (opt *Stop) Run(ctx app.Context) app.Error {
	opt.NoStyleArgs.Apply(&ctx)
	now := ctx.Now()
	date := opt.AtDate(now)
	time, err := opt.AtTime(now, ctx.Config())
	if err != nil {
		return err
	}
	// Only fall back to yesterday if no explicit date has been given.
	// Otherwise, it wouldn’t make sense to decrement the day.
	shouldTryYesterday := opt.WasAutomatic()
	yesterday := date.PlusDays(-1)
	return util.Reconcile(ctx, util.ReconcileOpts{OutputFileArgs: opt.OutputFileArgs, WarnArgs: opt.WarnArgs},
		[]reconciling.Creator{
			reconciling.NewReconcilerAtRecord(date),
			func() reconciling.Creator {
				if shouldTryYesterday {
					return reconciling.NewReconcilerAtRecord(yesterday)
				}
				return nil
			}(),
		},

		func(reconciler *reconciling.Reconciler) error {
			if shouldTryYesterday && reconciler.Record.Date().IsEqualTo(yesterday) {
				time, _ = time.Plus(klog.NewDuration(24, 0))
			}
			return reconciler.CloseOpenRange(time, opt.TimeFormat(ctx.Config()), opt.Summary)
		},
	)
}
