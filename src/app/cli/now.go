package cli

import (
	"errors"
	"fmt"
	. "klog"
	"klog/app"
	"klog/app/cli/lib"
	"klog/service"
	"os"
	"os/signal"
	"syscall"
	gotime "time"
)

type Now struct {
	lib.DiffArgs
	Follow bool `name:"follow" short:"f" help:"Keep shell open and follow changes"`
	lib.WarnArgs
	lib.InputFilesArgs
}

func (opt *Now) Run(ctx app.Context) error {
	h := func() error { return handle(opt, ctx) }
	if opt.Follow {
		return withRepeat(ctx, h)
	}
	return h()
}

func handle(opt *Now, ctx app.Context) error {
	now := ctx.Now()
	records, err := ctx.ReadInputs(opt.File...)
	if err != nil {
		return err
	}
	recents, err := getTodayOrYesterday(now, records)
	if err != nil {
		ctx.Print(err.Error())
		return nil
	}
	// Headline:
	label := "     Today"
	if !recents[0].Date().IsEqualTo(NewDateFromTime(now)) {
		label = " Yesterday"
	}
	ctx.Print("       " + label + "    " + "Overall\n")
	// Total:
	ctx.Print("Total  ")
	total, _ := service.HypotheticalTotal(now, recents...)
	grandTotal, _ := service.HypotheticalTotal(now, records...)
	ctx.Print(lib.Pad(10-len(total.ToString())) + lib.Styler.Duration(total, false))
	ctx.Print(lib.Pad(11-len(grandTotal.ToString())) + lib.Styler.Duration(grandTotal, false))
	ctx.Print("\n")
	if opt.Diff {
		// Should:
		ctx.Print("Should  ")
		shouldTotal := service.ShouldTotalSum(recents...)
		grandShouldTotal := service.ShouldTotalSum(records...)
		ctx.Print(lib.Pad(9-len(shouldTotal.ToString())) + lib.Styler.ShouldTotal(shouldTotal))
		ctx.Print(lib.Pad(11-len(grandShouldTotal.ToString())) + lib.Styler.ShouldTotal(grandShouldTotal))
		ctx.Print("\n")
		// Diff:
		ctx.Print("Diff    ")
		diff := service.Diff(shouldTotal, total)
		grandDiff := service.Diff(grandShouldTotal, grandTotal)
		ctx.Print(lib.Pad(9-len(diff.ToStringWithSign())) + lib.Styler.Duration(diff, true))
		ctx.Print(lib.Pad(11-len(grandDiff.ToStringWithSign())) + lib.Styler.Duration(grandDiff, true))
		ctx.Print("\n")
		// ETA:
		ctx.Print("E.T.A.  ")
		eta, _ := NewTimeFromTime(now).Add(NewDuration(0, 0).Minus(diff))
		if eta != nil {
			ctx.Print(lib.Pad(9-len(eta.ToString())) + lib.Styler.Time(eta))
		} else {
			ctx.Print(lib.Pad(9-3) + "???")
		}
		grandEta, _ := NewTimeFromTime(now).Add(NewDuration(0, 0).Minus(grandDiff))
		if grandEta != nil {
			ctx.Print(lib.Pad(11-len(grandEta.ToString())) + lib.Styler.Time(grandEta))
		} else {
			ctx.Print(lib.Pad(11-3) + "???")
		}
		ctx.Print("\n")
	}
	ctx.Print(opt.WarnArgs.ToString(now, records))
	return nil
}

func getTodayOrYesterday(now gotime.Time, records []Record) ([]Record, error) {
	rs := service.Sort(records, false)
	for i := 0; i <= 1; i++ {
		rs = service.Filter(records, service.FilterQry{
			Dates: []Date{NewDateFromTime(now).PlusDays(-i)},
		})
		if len(rs) > 0 {
			return rs, nil
		}
	}
	return nil, errors.New("No record found for today\n")
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
	ticker := gotime.NewTicker(1 * gotime.Second)
	defer ticker.Stop()
	i := 5 // seconds to display help text (how to exit)
	for ; true; <-ticker.C {
		ctx.Print(fmt.Sprintf("\033[H\033[J")) // Cursor reset
		err := fn()
		ctx.Print("\n")
		if i > 0 {
			ctx.Print("Press ^C to exit")
			i--
		}
		if err != nil {
			return err
		}
	}
	return nil
}
