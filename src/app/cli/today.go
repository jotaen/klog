package cli

import (
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

type Today struct {
	lib.DiffArgs
	lib.NowArgs
	Follow bool `name:"follow" short:"f" help:"Keep shell open and follow changes"`
	lib.WarnArgs
	lib.NoStyleArgs
	lib.InputFilesArgs
}

func (opt *Today) Help() string {
	return `Evaluates the total time, separately for today’s records and all other records.

When both --now and --diff are set, it also calculates the forecasted end-time at which the time goal will be reached.
(I.e. when the difference between should and actual time will be 0.)

If there are no records today, it falls back to yesterday.`
}

func (opt *Today) Run(ctx app.Context) error {
	opt.NoStyleArgs.Apply(&ctx)
	h := func() error { return handle(opt, ctx) }
	if opt.Follow {
		return withRepeat(ctx, h)
	}
	return h()
}

var (
	INDENT = "          "
	N_A    = "n/a"
	QQQ    = "???"
	COL_1  = 8
	COL_2  = 10
	COL_3  = 9
	COL_4  = 11
)

func handle(opt *Today, ctx app.Context) error {
	now := ctx.Now()
	records, err := ctx.ReadInputs(opt.File...)
	if err != nil {
		return err
	}

	currentRecords, otherRecords, isYesterday := splitIntoCurrentAndOther(now, records)
	hasCurrentRecords := len(currentRecords) > 0

	currentTotal, currentShouldTotal, currentDiff := opt.evaluate(now, currentRecords)
	currentEndTime, _ := NewTimeFromTime(now).Add(NewDuration(0, 0).Minus(currentDiff))

	otherTotal, otherShouldTotal, otherDiff := opt.evaluate(now, otherRecords)

	grandTotal := currentTotal.Plus(otherTotal)
	grandShouldTotal := NewShouldTotal(0, currentShouldTotal.Plus(otherShouldTotal).InMinutes())
	grandDiff := service.Diff(grandShouldTotal, grandTotal)
	grandEndTime, _ := NewTimeFromTime(now).Add(NewDuration(0, 0).Minus(grandDiff))

	// Headline:
	ctx.Print(INDENT + "   Total")
	if opt.Diff {
		ctx.Print("    Should     Diff")
		if opt.Now {
			ctx.Print("   End-Time")
		}
	}
	ctx.Print("\n")

	// Current:
	if isYesterday {
		ctx.Print("Yesterday ")
	} else {
		ctx.Print("Today     ")
	}
	if hasCurrentRecords {
		ctx.Print(lib.Pad(COL_1-len(currentTotal.ToString())) + ctx.Serialiser().Duration(currentTotal))
	} else {
		ctx.Print(lib.Pad(COL_1-len(N_A)) + N_A)
	}
	if opt.Diff {
		if hasCurrentRecords {
			ctx.Print(lib.Pad(COL_2-len(currentShouldTotal.ToString())) + ctx.Serialiser().ShouldTotal(currentShouldTotal))
			ctx.Print(lib.Pad(COL_3-len(currentDiff.ToStringWithSign())) + ctx.Serialiser().SignedDuration(currentDiff))
		} else {
			ctx.Print(lib.Pad(COL_2-len(N_A)) + N_A)
			ctx.Print(lib.Pad(COL_3-len(N_A)) + N_A)
		}
		if opt.Now {
			if hasCurrentRecords {
				if currentEndTime != nil {
					ctx.Print(lib.Pad(COL_4-len(currentEndTime.ToString())) + ctx.Serialiser().Time(currentEndTime))
				} else {
					ctx.Print(lib.Pad(COL_4-len(QQQ)) + QQQ)
				}
			} else {
				ctx.Print(lib.Pad(COL_4-len(N_A)) + N_A)
			}
		}
	}
	ctx.Print("\n")

	// Other:
	ctx.Print("Other     ")
	ctx.Print(lib.Pad(COL_1-len(otherTotal.ToString())) + ctx.Serialiser().Duration(otherTotal))
	if opt.Diff {
		ctx.Print(lib.Pad(COL_2-len(otherShouldTotal.ToString())) + ctx.Serialiser().ShouldTotal(otherShouldTotal))
		ctx.Print(lib.Pad(COL_3-len(otherDiff.ToStringWithSign())) + ctx.Serialiser().SignedDuration(otherDiff))
	}
	ctx.Print("\n")

	// Line:
	ctx.Print(INDENT + "========")
	if opt.Diff {
		ctx.Print("===================")
	}
	ctx.Print("\n")

	// GrandTotal:
	ctx.Print("All       ")
	ctx.Print(lib.Pad(COL_1-len(grandTotal.ToString())) + ctx.Serialiser().Duration(grandTotal))
	if opt.Diff {
		ctx.Print(lib.Pad(COL_2-len(grandShouldTotal.ToString())) + ctx.Serialiser().ShouldTotal(grandShouldTotal))
		ctx.Print(lib.Pad(COL_3-len(grandDiff.ToStringWithSign())) + ctx.Serialiser().SignedDuration(grandDiff))
		if opt.Now {
			if hasCurrentRecords {
				if grandEndTime != nil {
					ctx.Print(lib.Pad(COL_4-len(grandEndTime.ToString())) + ctx.Serialiser().Time(grandEndTime))
				} else {
					ctx.Print(lib.Pad(COL_4-len(QQQ)) + QQQ)
				}
			} else {
				ctx.Print(lib.Pad(COL_4-len(N_A)) + N_A)
			}
		}
	}
	ctx.Print("\n")

	ctx.Print(opt.WarnArgs.ToString(now, records))
	return nil
}

func (opt *Today) evaluate(now gotime.Time, records []Record) (Duration, Duration, Duration) {
	total, _ := func() (Duration, bool) {
		if opt.Now {
			return service.HypotheticalTotal(now, records...)
		}
		return service.Total(records...), false
	}()
	shouldTotal := service.ShouldTotalSum(records...)
	diff := service.Diff(shouldTotal, total)
	return total, shouldTotal, diff
}

func splitIntoCurrentAndOther(now gotime.Time, records []Record) ([]Record, []Record, bool) {
	var todaysRecords []Record
	var yesterdaysRecords []Record
	var otherRecords []Record
	today := NewDateFromTime(now)
	yesterday := today.PlusDays(-1)
	for _, r := range records {
		if r.Date().IsEqualTo(today) {
			todaysRecords = append(todaysRecords, r)
		} else if r.Date().IsEqualTo(yesterday) {
			yesterdaysRecords = append(yesterdaysRecords, r)
		} else {
			otherRecords = append(otherRecords, r)
		}
	}
	if len(todaysRecords) > 0 {
		return todaysRecords, append(otherRecords, yesterdaysRecords...), false
	}
	if len(yesterdaysRecords) > 0 {
		return yesterdaysRecords, otherRecords, true
	}
	return nil, otherRecords, false
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
