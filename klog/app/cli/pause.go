package cli

import (
	"fmt"
	"github.com/jotaen/klog/klog"
	"github.com/jotaen/klog/klog/app"
	"github.com/jotaen/klog/klog/app/cli/terminalformat"
	"github.com/jotaen/klog/klog/app/cli/util"
	"github.com/jotaen/klog/klog/parser"
	"github.com/jotaen/klog/klog/parser/reconciling"
	"strings"
	gotime "time"
)

type Pause struct {
	Summary klog.EntrySummary `name:"summary" short:"s" placeholder:"TEXT" help:"Summary text for the pause entry"`
	Extend  bool              `name:"extend" short:"e" help:"Extend latest pause, instead of adding a new pause entry"`
	util.OutputFileArgs
	util.NoStyleArgs
	util.WarnArgs
}

func (opt *Pause) Help() string {
	return `Creates a pause entry for a record with an open time range.
The command is blocking – it keeps updating the pause entry until the process is exited.
(The file will be written into once per minute.)
`
}

func (opt *Pause) Run(ctx app.Context) app.Error {
	opt.NoStyleArgs.Apply(&ctx)
	today := klog.NewDateFromGo(ctx.Now())
	doReconcile := func(reconcile reconciling.Reconcile) (*reconciling.Result, app.Error) {
		return ctx.ReconcileFile(
			opt.OutputFileArgs.File,
			[]reconciling.Creator{
				reconciling.NewReconcilerAtRecord(today),
				reconciling.NewReconcilerAtRecord(today.PlusDays(-1)),
			},
			reconcile,
		)
	}

	// Initial run:
	// Ensure that an open range exists, and set up the pause entry:
	// - Without `--extend`, append a new entry, including the summary
	// - With `--extend`, find a pause and append the summary
	lastResult, err := doReconcile(func(reconciler *reconciling.Reconciler) error {
		if opt.Extend {
			return reconciler.ExtendPause(klog.NewDuration(0, 0), opt.Summary)
		}
		return reconciler.AppendPause(opt.Summary)
	})
	if err != nil {
		return err
	}

	// Subsequent runs:
	// We don’t rely on the accumulated counter, because then it might also accumulate
	// imprecisions over time. Therefore, we always base the increment off the initial
	// start time. Also, if the computer is set to sleep, it should properly “recover”
	// afterwards.
	start := ctx.Now()
	minsCaptured := 0 // The amount of minutes that have already been written into the file.
	return util.WithRepeat(ctx.Print, 500*gotime.Millisecond, func(counter int64) app.Error {
		uncapturedIncrement := diffInMinutes(ctx.Now(), start) - minsCaptured
		ctx.Debug(func() {
			ctx.Print(fmt.Sprintf("Started: %s\n", start))
			ctx.Print(fmt.Sprintf("Now:     %s\n", ctx.Now()))
			ctx.Print(fmt.Sprintf("Incr.:   %d\n", uncapturedIncrement))
			ctx.Print("\n")
		})
		if uncapturedIncrement > 0 {
			lastResult, err = doReconcile(func(reconciler *reconciling.Reconciler) error {
				// Don’t add the summary here, as we already appended it in the initial run.
				return reconciler.ExtendPause(klog.NewDuration(0, -1*uncapturedIncrement), nil)
			})
			minsCaptured += uncapturedIncrement
			if err != nil {
				return err
			}
		}

		dots := strings.Repeat(".", int(counter%4))
		styler, serialiser := ctx.Serialise()
		ctx.Print("" +
			"Pausing for " +
			// Always print number in red, but without sign
			styler.Props(terminalformat.StyleProps{Color: terminalformat.RED}).Format(klog.NewDuration(0, minsCaptured).ToString()) +
			fmt.Sprintf("%-4s", dots) +
			"(since " +
			klog.NewTimeFromGo(start).ToString() +
			")\n")
		ctx.Print("\n" + parser.SerialiseRecords(serialiser, lastResult.Record).ToString() + "\n")
		if counter < 14 {
			// Display exit hint for a couple of seconds.
			ctx.Print("\n")
			ctx.Print("Press ^C to stop\n")
		}
		return nil
	})
}

// diffInMinutes computes the “wall-clock” difference between two times.
// Note, the built-in `Time.Sub` function computes the difference of the
// underlying monotonic time counter, which would yield incorrect results
// in case the monotonic timer was suspended, e.g. due to sleep.
func diffInMinutes(t1 gotime.Time, t2 gotime.Time) int {
	return int(t1.Unix()-t2.Unix()) / 60
}
