package lib

import (
	klog "github.com/jotaen/klog/src"
	"github.com/jotaen/klog/src/app"
	"github.com/jotaen/klog/src/parser"
	"github.com/jotaen/klog/src/parser/reconciling"
)

type ReconcileOpts struct {
	OutputFileArgs
	WarnArgs
}

func Reconcile(ctx app.Context, opts ReconcileOpts, creators []reconciling.Creator, reconcile reconciling.Reconcile) error {
	result, err := ctx.ReconcileFile(
		opts.OutputFileArgs.File,
		creators,
		reconcile,
	)
	if err != nil {
		return err
	}
	ctx.Print("\n" + ctx.Serialiser().SerialiseRecords(result.Record) + "\n")
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
