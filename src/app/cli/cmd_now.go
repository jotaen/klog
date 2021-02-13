package cli

import (
	"errors"
	"fmt"
	. "klog"
	"klog/app"
	"klog/service"
	"os"
	"os/signal"
	"syscall"
	"time"
)

type Now struct {
	DiffArg
	Follow bool `name:"follow" help:"Keep shell open and follow changes"`
	WarnArgs
	InputFilesArgs
}

func (args *Now) Run(ctx app.Context) error {
	handle := func() error {
		records, err := ctx.RetrieveRecords(args.File...)
		if err != nil {
			return err
		}
		now := time.Now()
		recent, err := func() (Record, error) {
			rs := service.Query(records, service.Opts{
				Dates: []Date{NewDateFromTime(now), NewDateFromTime(now).PlusDays(-1)},
				Sort:  "DESC",
			})
			if len(rs) == 0 {
				return nil, errors.New("No record found for today\n")
			}
			return rs[0], nil
		}()
		if err != nil {
			ctx.Print("Error: " + err.Error())
			return nil
		}
		// Headline:
		label := "     Today"
		if !recent.Date().IsEqualTo(NewDateFromTime(now)) {
			label = " Yesterday"
		}
		ctx.Print("       " + label + "    " + "Overall\n")
		// Total:
		ctx.Print("Total  ")
		total, _ := service.HypotheticalTotal(now, recent)
		grandTotal, _ := service.HypotheticalTotal(now, records...)
		ctx.Print(pad(10-len(total.ToString())) + styler.Duration(total, false))
		ctx.Print(pad(11-len(grandTotal.ToString())) + styler.Duration(grandTotal, false))
		ctx.Print("\n")
		if args.Diff {
			// Should:
			ctx.Print("Should  ")
			shouldTotal := service.ShouldTotalSum(recent)
			grandShouldTotal := service.ShouldTotalSum(records...)
			ctx.Print(pad(9-len(shouldTotal.ToString())) + styler.ShouldTotal(shouldTotal))
			ctx.Print(pad(11-len(grandShouldTotal.ToString())) + styler.ShouldTotal(grandShouldTotal))
			ctx.Print("\n")
			// Diff:
			ctx.Print("Diff    ")
			diff := total.Minus(shouldTotal)
			grandDiff := grandTotal.Minus(grandShouldTotal)
			ctx.Print(pad(9-len(diff.ToStringWithSign())) + styler.Duration(diff, true))
			ctx.Print(pad(11-len(grandDiff.ToStringWithSign())) + styler.Duration(grandDiff, true))
			ctx.Print("\n")
			// ETA:
			ctx.Print("E.T.A.  ")
			eta, _ := NewTimeFromTime(now).Add(NewDuration(0, 0).Minus(diff))
			if eta != nil {
				ctx.Print(pad(9-len(eta.ToString())) + styler.Time(eta))
			} else {
				ctx.Print(pad(9-3) + "???")
			}
			grandEta, _ := NewTimeFromTime(now).Add(NewDuration(0, 0).Minus(grandDiff))
			if grandEta != nil {
				ctx.Print(pad(11-len(grandEta.ToString())) + styler.Time(grandEta))
			} else {
				ctx.Print(pad(11-3) + "???")
			}
			ctx.Print("\n")
		}
		args.printWarnings(ctx, records)
		return nil
	}
	if args.Follow {
		return withRepeat(ctx, handle)
	}
	return handle()
}

func withRepeat(ctx app.Context, fn func() error) error {
	// Handle ^C gracefully, as it’s the only way to exit
	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		os.Exit(0)
		return
	}()

	// Call handler function repetitively
	ctx.Print("\033[2J") // Initial screen clearing
	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()
	for ; true; <-ticker.C {
		ctx.Print(fmt.Sprintf("\033[H\033[J")) // Cursor reset
		err := fn()
		ctx.Print("\n")
		if err != nil {
			return err
		}
	}
	return nil
}
