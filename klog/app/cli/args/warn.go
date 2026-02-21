package args

import (
	"github.com/jotaen/klog/klog"
	"github.com/jotaen/klog/klog/app"
	"github.com/jotaen/klog/klog/app/cli/prettify"
	"github.com/jotaen/klog/klog/service"
)

type WarnArgs struct {
	NoWarn bool `name:"no-warn" help:"Suppress warnings about potential mistakes or logical errors."`
}

func (args *WarnArgs) PrintWarnings(ctx app.Context, records []klog.Record, additionalWarnings []service.UsageWarning) {
	styler, _ := ctx.Serialise()
	warnings := args.GatherWarnings(ctx, records, additionalWarnings)
	for _, w := range warnings {
		ctx.Print(prettify.PrettifyWarning(w, styler))
	}
}

func (args *WarnArgs) GatherWarnings(ctx app.Context, records []klog.Record, additionalWarnings []service.UsageWarning) []string {
	if args.NoWarn {
		return nil
	}
	disabledCheckers := ctx.Config().NoWarnings.UnwrapOr(service.NewDisabledCheckers())
	warnings := service.CheckForWarnings(ctx.Now(), records, disabledCheckers)
	for _, warn := range additionalWarnings {
		if warn != (service.UsageWarning{}) && !disabledCheckers[warn.Name] {
			warnings = append(warnings, warn.Message)
		}
	}
	return warnings
}
