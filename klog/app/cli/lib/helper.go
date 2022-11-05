package lib

import (
	"github.com/jotaen/klog/klog/app"
	"github.com/jotaen/klog/klog/parser"
	"github.com/jotaen/klog/klog/parser/reconciling"
)

type ReconcileOpts struct {
	OutputFileArgs
	WarnArgs
}

func Reconcile(ctx app.Context, opts ReconcileOpts, creators []reconciling.Creator, reconcile reconciling.Reconcile) error {
	result, err := ctx.ReconcileFile(
		true,
		opts.OutputFileArgs.File,
		creators,
		reconcile,
	)
	if err != nil {
		return err
	}
	ctx.Print("\n" + parser.SerialiseRecords(ctx.Serialiser(), result.Record).ToString() + "\n")
	opts.WarnArgs.PrintWarnings(ctx, result.AllRecords)
	return nil
}
