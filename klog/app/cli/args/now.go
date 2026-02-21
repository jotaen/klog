package args

import (
	gotime "time"

	"github.com/jotaen/klog/klog"
	"github.com/jotaen/klog/klog/app"
	"github.com/jotaen/klog/klog/service"
)

type NowArgs struct {
	Now          bool `name:"now" short:"n" help:"Assume open ranges to be closed at this moment."`
	hadOpenRange bool // Field only for internal use
}

func (args *NowArgs) ApplyNow(reference gotime.Time, rs ...klog.Record) app.Error {
	if args.Now {
		hasClosedAnyRange, err := service.CloseOpenRanges(reference, rs...)
		if err != nil {
			return app.NewErrorWithCode(
				app.LOGICAL_ERROR,
				"Cannot apply --now flag",
				"There are records with uncloseable time ranges",
				err,
			)
		}
		args.hadOpenRange = hasClosedAnyRange
		return nil
	}
	return nil
}

func (args *NowArgs) HadOpenRange() bool {
	return args.hadOpenRange
}

func (args *NowArgs) GetNowWarnings() []string {
	if args.Now && !args.hadOpenRange {
		return []string{"You specified --now, but there was no open-ended time range"}
	}
	return nil
}
