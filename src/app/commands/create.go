package commands

import (
	. "github.com/jotaen/klog/src"
	"github.com/jotaen/klog/src/app"
	"github.com/jotaen/klog/src/app/lib"
	"github.com/jotaen/klog/src/parser"
	"github.com/jotaen/klog/src/parser/reconciling"
)

type Create struct {
	ShouldTotal ShouldTotal `name:"should" help:"The should-total of the record"`
	lib.AtDateArgs
	lib.NoStyleArgs
	lib.OutputFileArgs
}

func (opt *Create) Help() string {
	return `The new record is inserted into the file at the chronologically correct position.
(Assuming that the records are sorted from oldest to latest.)`
}

func (opt *Create) Run(ctx app.Context) error {
	opt.NoStyleArgs.Apply(&ctx)
	date, _ := opt.AtDate(ctx.Now())
	return ctx.ReconcileFile(
		opt.OutputFileArgs.File,

		[]reconciling.Creator{
			func(parsedRecords []parser.ParsedRecord) *reconciling.Reconciler {
				return reconciling.NewReconcilerAtNewRecord(parsedRecords, date, opt.ShouldTotal)
			},
		},

		func(reconciler *reconciling.Reconciler) (*reconciling.Result, error) {
			return reconciler.MakeResult()
		},
	)
}
