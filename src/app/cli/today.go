package cli

import (
	. "github.com/jotaen/klog/src"
	"github.com/jotaen/klog/src/app"
	"github.com/jotaen/klog/src/app/cli/lib"
	"github.com/jotaen/klog/src/app/cli/lib/terminalformat"
	"github.com/jotaen/klog/src/service"
	gotime "time"
)

type Today struct {
	lib.DiffArgs
	lib.NowArgs
	Follow bool `name:"follow" short:"f" help:"Keep shell open and follow changes"`
	lib.DecimalArgs
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
	opt.DecimalArgs.Apply(&ctx)
	opt.NoStyleArgs.Apply(&ctx)
	if opt.Follow {
		return lib.WithRepeat(ctx.Print, 1*gotime.Second, func(counter int64) error {
			err := handle(opt, ctx)
			if counter < 7 {
				// Display exit hint for a couple of seconds.
				ctx.Print("\n")
				ctx.Print("Press ^C to exit")
				ctx.Print("\n")
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

func handle(opt *Today, ctx app.Context) error {
	records, err := ctx.ReadInputs(opt.File...)
	if err != nil {
		return err
	}
	now := ctx.Now()
	records, nErr := opt.ApplyNow(now, records...)
	if nErr != nil {
		return nil
	}

	currentRecords, otherRecords, isYesterday := splitIntoCurrentAndOther(now, records)
	hasCurrentRecords := len(currentRecords) > 0

	currentTotal, currentShouldTotal, currentDiff := opt.evaluate(currentRecords)
	currentEndTime, _ := NewTimeFromGo(now).Plus(NewDuration(0, 0).Minus(currentDiff))

	otherTotal, otherShouldTotal, otherDiff := opt.evaluate(otherRecords)

	grandTotal := currentTotal.Plus(otherTotal)
	grandShouldTotal := NewShouldTotal(0, currentShouldTotal.Plus(otherShouldTotal).InMinutes())
	grandDiff := service.Diff(grandShouldTotal, grandTotal)
	grandEndTime, _ := NewTimeFromGo(now).Plus(NewDuration(0, 0).Minus(grandDiff))

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
	table := terminalformat.NewTable(numberOfColumns, " ")

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
		table.CellR(ctx.Serialiser().Duration(currentTotal))
	} else {
		table.CellR(N_A)
	}
	if opt.Diff {
		if hasCurrentRecords {
			table.
				CellR(ctx.Serialiser().ShouldTotal(currentShouldTotal)).
				CellR(ctx.Serialiser().SignedDuration(currentDiff))
		} else {
			table.CellR(N_A).CellR(N_A)
		}
		if opt.Now {
			if hasCurrentRecords {
				if currentEndTime != nil {
					table.CellR(ctx.Serialiser().Time(currentEndTime))
				} else {
					table.CellR(QQQ)
				}
			} else {
				table.CellR(N_A)
			}
		}
	}

	// Other:
	table.CellL("Other").CellR(ctx.Serialiser().Duration(otherTotal))
	if opt.Diff {
		table.
			CellR(ctx.Serialiser().ShouldTotal(otherShouldTotal)).
			CellR(ctx.Serialiser().SignedDuration(otherDiff))
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
	table.CellL("All").CellR(ctx.Serialiser().Duration(grandTotal))
	if opt.Diff {
		table.
			CellR(ctx.Serialiser().ShouldTotal(grandShouldTotal)).
			CellR(ctx.Serialiser().SignedDuration(grandDiff))
		if opt.Now {
			if hasCurrentRecords {
				if grandEndTime != nil {
					table.CellR(ctx.Serialiser().Time(grandEndTime))
				} else {
					table.CellR(QQQ)
				}
			} else {
				table.CellR(N_A)
			}
		}
	}
	table.Collect(ctx.Print)
	opt.WarnArgs.PrintWarnings(ctx, records)
	return nil
}

func (opt *Today) evaluate(records []Record) (Duration, Duration, Duration) {
	total := service.Total(records...)
	shouldTotal := service.ShouldTotalSum(records...)
	diff := service.Diff(shouldTotal, total)
	return total, shouldTotal, diff
}

func splitIntoCurrentAndOther(now gotime.Time, records []Record) ([]Record, []Record, bool) {
	var todaysRecords []Record
	var yesterdaysRecords []Record
	var otherRecords []Record
	today := NewDateFromGo(now)
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
