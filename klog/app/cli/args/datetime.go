package args

import (
	gotime "time"

	"github.com/jotaen/klog/klog"
	"github.com/jotaen/klog/klog/app"
	"github.com/jotaen/klog/klog/parser/reconciling"
	"github.com/jotaen/klog/klog/service"
)

type AtDateArgs struct {
	Date      klog.Date `name:"date" placeholder:"DATE" short:"d" help:"The date of the record."`
	Today     bool      `name:"today" help:"Use today’s date."`
	Yesterday bool      `name:"yesterday" help:"Use yesterday’s date."`
	Tomorrow  bool      `name:"tomorrow" help:"Use tomorrow’s date."`
}

func (args *AtDateArgs) AtDate(now gotime.Time) klog.Date {
	if args.Date != nil {
		return args.Date
	}
	today := klog.NewDateFromGo(now) // That’s effectively/implicitly `--today`
	if args.Yesterday {
		return today.PlusDays(-1)
	} else if args.Tomorrow {
		return today.PlusDays(1)
	}
	return today
}

func (args *AtDateArgs) DateFormat(config app.Config) reconciling.ReformatDirective[klog.DateFormat] {
	if args.Date != nil {
		return reconciling.NoReformat[klog.DateFormat]()
	}
	fd := reconciling.ReformatAutoStyle[klog.DateFormat]()
	config.DateUseDashes.Unwrap(func(x bool) {
		fd = reconciling.ReformatExplicitly(klog.DateFormat{UseDashes: x})
	})
	return fd
}

type AtDateAndTimeArgs struct {
	Round service.Rounding `name:"round" placeholder:"ROUNDING" short:"r" help:"Round time to nearest multiple number. ROUNDING can be one of '5m', '10m', '12m', '15m', '20m', '30m' or '60m' / '1h'."`
	AtDateArgs
	Time klog.Time `name:"time" placeholder:"TIME" short:"t" help:"Specify the time (defaults to now). TIME can be given in the 24h or 12h notation, e.g. '13:00' or '1:00pm'."`
}

func (args *AtDateAndTimeArgs) AtTime(now gotime.Time, config app.Config) (klog.Time, app.Error) {
	if args.Time != nil {
		return args.Time, nil
	}
	date := args.AtDate(now)
	today := klog.NewDateFromGo(now)
	time := klog.NewTimeFromGo(now)
	if args.Round != nil {
		time = service.RoundToNearest(time, args.Round)
	} else {
		config.DefaultRounding.Unwrap(func(r service.Rounding) {
			time = service.RoundToNearest(time, r)
		})
	}
	if today.IsEqualTo(date) {
		return time, nil
	} else if today.PlusDays(-1).IsEqualTo(date) {
		shiftedTime, _ := time.Plus(klog.NewDuration(24, 0))
		return shiftedTime, nil
	} else if today.PlusDays(1).IsEqualTo(date) {
		shiftedTime, _ := time.Plus(klog.NewDuration(-24, 0))
		return shiftedTime, nil
	}
	return nil, app.NewErrorWithCode(
		app.LOGICAL_ERROR,
		"Missing time parameter",
		"Please specify a time value for dates in the past",
		nil,
	)
}

func (args *AtDateAndTimeArgs) TimeFormat(config app.Config) reconciling.ReformatDirective[klog.TimeFormat] {
	if args.Time != nil {
		return reconciling.NoReformat[klog.TimeFormat]()
	}
	fd := reconciling.ReformatAutoStyle[klog.TimeFormat]()
	config.TimeUse24HourClock.Unwrap(func(x bool) {
		fd = reconciling.ReformatExplicitly(klog.TimeFormat{Use24HourClock: x})
	})
	return fd
}

func (args *AtDateAndTimeArgs) WasAutomatic() bool {
	return args.Date == nil && args.Time == nil
}
