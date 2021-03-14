package cli

import (
	"errors"
	. "klog"
	"klog/app"
	"klog/app/cli/lib"
	"klog/parser"
)

type Track struct {
	lib.AtDateArgs
	Entry string `arg required help:"The new entry to add (requires quoting if contains multiple words)"`
	lib.OutputFileArgs
}

func (opt *Track) Run(ctx app.Context) error {
	date := opt.AtDate(ctx.Now())
	return reconcile(
		opt.OutputFileArgs,
		ctx,
		func(pr *parser.ParseResult) (*parser.Reconciler, error) {
			return parser.NewRecordReconciler(pr,
				errors.New("No record at date "+date.ToString()),
				func(r Record) bool { return r.Date().IsEqualTo(date) })
		},
		func(r *parser.Reconciler) (Record, string, error) {
			return r.AppendEntry(
				func(r Record) string { return opt.Entry },
			)
		},
	)
}

func reconcile(
	fileArgs lib.OutputFileArgs,
	ctx app.Context,
	reconcilerSupplier func(pr *parser.ParseResult) (*parser.Reconciler, error),
	handle func(reconciler *parser.Reconciler) (Record, string, error),
) error {
	targetFile, err := fileArgs.OutputFile(ctx)
	if err != nil {
		return err
	}
	pr, err := ctx.ReadFileInput(targetFile)
	if err != nil {
		return err
	}
	reconciler, err := reconcilerSupplier(pr)
	if reconciler == nil {
		return err
	}
	record, contents, err := handle(reconciler)
	if err != nil {
		return err
	}
	err = ctx.WriteFile(targetFile, contents)
	if err != nil {
		return err
	}
	ctx.Print("\n" + parser.SerialiseRecords(&lib.Styler, record) + "\n")
	return nil
}
