package cli

import (
	"github.com/jotaen/klog/klog"
	"github.com/jotaen/klog/klog/app"
	terminalformat2 "github.com/jotaen/klog/klog/app/cli/terminalformat"
	"github.com/jotaen/klog/klog/app/cli/util"
	"github.com/jotaen/klog/klog/service"
	gotime "time"
)

type Today struct {
	util.DiffArgs
	util.NowArgs
	Follow bool `name:"follow" short:"f" help:"Keep shell open and follow changes"`
	util.DecimalArgs
	util.WarnArgs
	util.NoStyleArgs
	util.InputFilesArgs
}

func (opt *Today) Help() string {
	return `Convenience command to “check in” on the current day.
It evaluates the total time separately for today’s records and all other records.

When both --now and --diff are set, it also calculates the forecasted end-time at which the time goal will be reached.
(I.e. when the difference between should and actual time will be 0.)

If there are no records today, it falls back to yesterday.`
}

func (opt *Today) Run(ctx app.Context) app.Error {
	opt.DecimalArgs.Apply(&ctx)
	opt.NoStyleArgs.Apply(&ctx)
	if opt.Follow {
		return util.WithRepeat(ctx.Print, 1*gotime.Second, func(counter int64) app.Error {
			err := handle(opt, ctx)
			if counter < 7 {
				// Display exit hint for a couple of seconds.
				ctx.Print("\nPress ^C to exit\n")
			}
			return err
		})
	}
	return handle(opt, ctx)
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

func handle(opt *Today, ctx app.Context) app.Error {
	records, err := ctx.ReadInputs(opt.File...)
	if err != nil {
		return err
	}
	now := ctx.Now()
	nErr := opt.ApplyNow(now, records...)
	if nErr != nil {
		return nErr
	}
	styler, serialiser := ctx.Serialise()

	currentRecords, otherRecords, isYesterday := splitIntoCurrentAndOther(now, records)
	hasCurrentRecords := len(currentRecords) > 0

	currentTotal, currentShouldTotal, currentDiff := opt.evaluate(currentRecords)
	currentEndTime, _ := klog.NewTimeFromGo(now).Plus(klog.NewDuration(0, 0).Minus(currentDiff))

	otherTotal, otherShouldTotal, otherDiff := opt.evaluate(otherRecords)

	grandTotal := currentTotal.Plus(otherTotal)
	grandShouldTotal := klog.NewShouldTotal(0, currentShouldTotal.Plus(otherShouldTotal).InMinutes())
	grandDiff := service.Diff(grandShouldTotal, grandTotal)
	grandEndTime, _ := klog.NewTimeFromGo(now).Plus(klog.NewDuration(0, 0).Minus(grandDiff))

	numberOfValueColumns := func() int {
		if opt.Diff {
			if opt.Now {
				return 4
			}
			return 3
		}
		return 1
	}()
	numberOfColumns := 1 + numberOfValueColumns
	table := terminalformat2.NewTable(numberOfColumns, " ")

	// Headline:
	table.
		CellL("         ").
		CellR("   Total")
	if opt.Diff {
		table.CellR("   Should").CellR("    Diff")
		if opt.Now {
			table.CellR("  End-Time")
		}
	}

	// Current:
	if isYesterday {
		table.CellL("Yesterday")
	} else {
		table.CellL("Today")
	}
	if hasCurrentRecords {
		table.CellR(serialiser.Duration(currentTotal))
	} else {
		table.CellR(N_A)
	}
	if opt.Diff {
		if hasCurrentRecords {
			table.
				CellR(serialiser.ShouldTotal(currentShouldTotal)).
				CellR(serialiser.SignedDuration(currentDiff))
		} else {
			table.CellR(N_A).CellR(N_A)
		}
		if opt.Now {
			if hasCurrentRecords {
				if currentEndTime != nil {
					if opt.HadOpenRange() {
						table.CellR(serialiser.Time(currentEndTime))
					} else {
						table.CellR(
							styler.Props(terminalformat2.StyleProps{Color: terminalformat2.SUBDUED}).
								Format("(" + currentEndTime.ToString() + ")"))
					}
				} else {
					table.CellR(QQQ)
				}
			} else {
				table.CellR(N_A)
			}
		}
	}

	// Other:
	table.CellL("Other").CellR(serialiser.Duration(otherTotal))
	if opt.Diff {
		table.
			CellR(serialiser.ShouldTotal(otherShouldTotal)).
			CellR(serialiser.SignedDuration(otherDiff))
		if opt.Now {
			table.Skip(1)
		}
	}

	// Line:
	table.Skip(1).Fill("=")
	if opt.Diff {
		table.Fill("=").Fill("=")
		if opt.Now {
			table.Skip(1)
		}
	}

	// GrandTotal:
	table.CellL("All").CellR(serialiser.Duration(grandTotal))
	if opt.Diff {
		table.
			CellR(serialiser.ShouldTotal(grandShouldTotal)).
			CellR(serialiser.SignedDuration(grandDiff))
		if opt.Now {
			if hasCurrentRecords {
				if grandEndTime != nil {
					if opt.HadOpenRange() {
						table.CellR(serialiser.Time(grandEndTime))
					} else {
						table.CellR(
							styler.Props(terminalformat2.StyleProps{Color: terminalformat2.SUBDUED}).
								Format("(" + grandEndTime.ToString() + ")"))
					}
				} else {
					table.CellR(QQQ)
				}
			} else {
				table.CellR(N_A)
			}
		}
	}
	table.Collect(ctx.Print)
	opt.WarnArgs.PrintWarnings(ctx, records, opt.GetNowWarnings())
	return nil
}

func (opt *Today) evaluate(records []klog.Record) (klog.Duration, klog.Duration, klog.Duration) {
	total := service.Total(records...)
	shouldTotal := service.ShouldTotalSum(records...)
	diff := service.Diff(shouldTotal, total)
	return total, shouldTotal, diff
}

func splitIntoCurrentAndOther(now gotime.Time, records []klog.Record) ([]klog.Record, []klog.Record, bool) {
	var todaysRecords []klog.Record
	var yesterdaysRecords []klog.Record
	var otherRecords []klog.Record
	today := klog.NewDateFromGo(now)
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
