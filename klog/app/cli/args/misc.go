package args

import (
	"strings"

	"github.com/jotaen/klog/klog"
	"github.com/jotaen/klog/klog/app"
	"github.com/jotaen/klog/klog/service"
	tf "github.com/jotaen/klog/lib/terminalformat"
)

type DiffArgs struct {
	Diff bool `name:"diff" short:"d" help:"Show difference between actual and should-total time."`
}

// GetWarning returns a warning if the user applied entry-level filtering (partial
// records) *and* requested to compute the should-total diff, as that may yield
// nonsensical results.
func (args *DiffArgs) GetWarning(filterArgs FilterArgs) service.UsageWarning {
	if args.Diff && filterArgs.hasPartialRecordsWithShouldTotal {
		return service.EntryFilteredDiffWarning
	}
	return service.UsageWarning{}
}

type NoStyleArgs struct {
	NoStyle bool `name:"no-style" help:"Do not style or colour the values."`
}

func (args *NoStyleArgs) Apply(ctx *app.Context) {
	if args.NoStyle {
		(*ctx).ConfigureSerialisation(func(styler tf.Styler, decimalDuration bool) (tf.Styler, bool) {
			return tf.NewStyler(tf.COLOUR_THEME_NO_COLOUR), decimalDuration
		})
	}
}

type QuietArgs struct {
	Quiet bool `name:"quiet" help:"Output parseable data without descriptive text."`
}

type SortArgs struct {
	Sort string `name:"sort" placeholder:"ORDER" help:"Sort output by date. ORDER can be 'asc' or 'desc'." enum:"asc,desc,ASC,DESC," default:""`
}

func (args *SortArgs) ApplySort(rs []klog.Record) []klog.Record {
	if args.Sort == "" {
		return rs
	}
	startWithOldest := false
	if strings.ToLower(args.Sort) == "asc" {
		startWithOldest = true
	}
	return service.Sort(rs, startWithOldest)
}

type DecimalArgs struct {
	Decimal bool `name:"decimal" help:"Display totals as decimal values (in minutes)."`
}

func (args *DecimalArgs) Apply(ctx *app.Context) {
	if args.Decimal {
		(*ctx).ConfigureSerialisation(func(styler tf.Styler, decimalDuration bool) (tf.Styler, bool) {
			return styler, true
		})
	}
}
