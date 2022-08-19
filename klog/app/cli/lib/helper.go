package lib

import (
	klog "github.com/jotaen/klog/klog"
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
	ctx.Print("\n" + parser.SerialiseRecords(ctx.Serialiser(), result.Record) + "\n")
	opts.WarnArgs.PrintWarnings(ctx, ToRecords(result.AllRecords))
	return nil
}

func ToRecords(prs []parser.ParsedRecord) []klog.Record {
	result := make([]klog.Record, len(prs))
	for i, r := range prs {
		result[i] = r
	}
	return result
}
