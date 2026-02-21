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
	if args.NoWarn {
		return
	}
	disabledCheckers := ctx.Config().NoWarnings.UnwrapOr(service.NewDisabledCheckers())
	for _, warn := range additionalWarnings {
		if warn != (service.UsageWarning{}) && !disabledCheckers[warn.Name] {
			ctx.Print(prettify.PrettifyGeneralWarning(warn.Message, styler))
		}
	}
	service.CheckForWarnings(func(w service.Warning) {
		ctx.Print(prettify.PrettifyWarning(w, styler))
	}, ctx.Now(), records, disabledCheckers)
}
