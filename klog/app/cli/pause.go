package cli

import (
	"github.com/jotaen/klog/klog"
	"github.com/jotaen/klog/klog/app"
	"github.com/jotaen/klog/klog/app/cli/lib"
	"github.com/jotaen/klog/klog/parser/reconciling"
	"strings"
	gotime "time"
)

type Pause struct {
	Summary klog.EntrySummary `name:"summary" short:"s" placeholder:"TEXT" help:"Summary text for the pause entry"`
	Extend  bool              `name:"extend" short:"e" help:"Extend latest pause, instead of adding a new pause entry"`
	lib.OutputFileArgs
	lib.NoStyleArgs
	lib.WarnArgs
}

func (opt *Pause) Help() string {
	return `Creates a pause entry for a record with an open time range.
The command is blocking – it keeps updating the pause entry until the process is exited.
(The file will be written into once per minute.)
`
}

func (opt *Pause) Run(ctx app.Context) app.Error {
	opt.NoStyleArgs.Apply(&ctx)
	today := klog.NewDateFromGo(gotime.Now())
	doReconcile := func(reconcile reconciling.Reconcile) app.Error {
		_, err := ctx.ReconcileFile(
			opt.OutputFileArgs.File,
			[]reconciling.Creator{
				reconciling.NewReconcilerAtRecord(today),
				reconciling.NewReconcilerAtRecord(today.PlusDays(-1)),
			},
			reconcile,
		)
		return err
	}

	// Initial run:
	// Ensure that an open range exists, and set up the pause entry:
	// - Without `--extend`, append a new entry, including the summary
	// - With `--extend`, find a pause and append the summary
	err := doReconcile(func(reconciler *reconciling.Reconciler) (*reconciling.Result, error) {
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
	// start time.
	start := gotime.Now()
	minsCaptured := 0 // The amount of minutes that have already been written into the file.
	return lib.WithRepeat(ctx.Print, 500*gotime.Millisecond, func(counter int64) app.Error {
		dots := strings.Repeat(".", int(counter%4))
		ctx.Print("" +
			"Pausing for " +
			// Always print number in red, but without sign
			ctx.Serialiser().Format(lib.Red, klog.NewDuration(0, minsCaptured).ToString()) +
			dots + "\n" +
			"(since " +
			klog.NewTimeFromGo(start).ToString() +
			")\n")
		if counter < 14 {
			// Display exit hint for a couple of seconds.
			ctx.Print("\n")
			ctx.Print("Press ^C to stop\n")
		}

		diffMins := int(gotime.Time.Sub(gotime.Now(), start).Minutes())
		increment := diffMins - minsCaptured
		if increment > 0 {
			minsCaptured += increment
			err := doReconcile(func(reconciler *reconciling.Reconciler) (*reconciling.Result, error) {
				// Don’t add the summary, as we already appended it in the initial run.
				return reconciler.ExtendPause(klog.NewDuration(0, -1*increment), nil)
			})
			if err != nil {
				return err
			}
			ctx.Debug(func() {
				ctx.Print("File saved.\n")
			})
		}
		return nil
	})
}
