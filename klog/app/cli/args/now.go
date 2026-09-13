package args

import (
	gotime "time"

	"github.com/jotaen/klog/klog"
	"github.com/jotaen/klog/klog/app"
	"github.com/jotaen/klog/klog/service"
)

type NowArgs struct {
	Now           bool                `name:"now" short:"n" help:"Assume open ranges to be closed at this moment."`
	closedEntries map[klog.Record]int // Field only for internal use
}

func (args *NowArgs) ApplyNow(reference gotime.Time, rs ...klog.Record) (map[klog.Record]int, app.Error) {
	if args.Now {
		closedEntries, err := service.CloseOpenRanges(reference, rs...)
		if err != nil {
			return nil, app.NewErrorWithCode(
				app.LOGICAL_ERROR,
				"Cannot apply --now flag",
				"There are records with uncloseable time ranges",
				err,
			)
		}
		args.closedEntries = closedEntries
		return closedEntries, nil
	}
	return nil, nil
}

func (args *NowArgs) HadOpenRange() bool {
	return len(args.closedEntries) > 0
}

// GetWarning warns the user that they specified the --now flag but there actually
// weren’t any closable ranges in the data.
func (args *NowArgs) GetWarning() service.UsageWarning {
	if args.Now && !args.HadOpenRange() {
		return service.PointlessNowWarning
	}
	return service.UsageWarning{}
}
