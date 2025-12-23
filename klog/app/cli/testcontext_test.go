package cli

import (
	gotime "time"

	"github.com/jotaen/klog/klog"
	"github.com/jotaen/klog/klog/app"
	"github.com/jotaen/klog/klog/app/cli/command"
	"github.com/jotaen/klog/klog/parser"
	"github.com/jotaen/klog/klog/parser/reconciling"
	"github.com/jotaen/klog/klog/parser/txt"
	tf "github.com/jotaen/klog/lib/terminalformat"
)

func NewTestingContext() TestingContext {
	bc := app.NewEmptyBookmarksCollection()
	config := app.NewDefaultConfig(tf.COLOUR_THEME_NO_COLOUR)
	styler := tf.NewStyler(tf.COLOUR_THEME_NO_COLOUR)
	return TestingContext{
		State: State{
			printBuffer:         "",
			writtenFileContents: "",
		},
		now:            gotime.Now(),
		records:        nil,
		blocks:         nil,
		styler:         styler,
		serialiser:     app.NewSerialiser(styler, false),
		bookmarks:      bc,
		editorsAuto:    nil,
		editorExplicit: "",
		fileExplorers:  nil,
		execute: func(_ command.Command) app.Error {
			return nil
		},
		config: &config,
	}
}

func (ctx TestingContext) _SetRecords(recordsText string) TestingContext {
	records, blocks, err := parser.NewSerialParser().Parse(recordsText)
	if err != nil {
		panic("Invalid records")
	}
	ctx.records = records
	ctx.blocks = blocks
	return ctx
}

func (ctx TestingContext) _SetNow(Y int, M int, D int, h int, m int) TestingContext {
	ctx.now = gotime.Date(Y, gotime.Month(M), D, h, m, 0, 0, gotime.UTC)
	return ctx
}

func (ctx TestingContext) _SetEditors(auto []command.Command, explicit string) TestingContext {
	ctx.editorsAuto = auto
	ctx.editorExplicit = explicit
	return ctx
}

func (ctx TestingContext) _SetFileExplorers(cs []command.Command) TestingContext {
	ctx.fileExplorers = cs
	return ctx
}

func (ctx TestingContext) _SetFileConfig(configFile string) TestingContext {
	cfg, err := app.NewConfig(1, func(_ string) string { return "" }, configFile)
	if err != nil {
		panic(err)
	}
	ctx.config = &cfg
	return ctx
}

func (ctx TestingContext) _SetExecute(execute func(command.Command) app.Error) TestingContext {
	ctx.execute = execute
	return ctx
}

func (ctx TestingContext) _Run(cmd func(app.Context) app.Error) (State, app.Error) {
	cmdErr := cmd(&ctx)
	out := ctx.printBuffer
	if len(out) > 0 && out[0] != '\n' {
		out = "\n" + out
	}
	return State{out, ctx.writtenFileContents}, cmdErr
}

type State struct {
	printBuffer         string
	writtenFileContents string
}

type TestingContext struct {
	State
	now            gotime.Time
	records        []klog.Record
	blocks         []txt.Block
	styler         tf.Styler
	serialiser     app.TextSerialiser
	bookmarks      app.BookmarksCollection
	editorsAuto    []command.Command
	editorExplicit string
	fileExplorers  []command.Command
	execute        func(command.Command) app.Error
	config         *app.Config
}

func (ctx *TestingContext) Print(s string) {
	ctx.printBuffer += s
}

func (ctx *TestingContext) ReadLine() (string, app.Error) {
	return "", nil
}

func (ctx *TestingContext) HomeFolder() string {
	return "~"
}

func (ctx *TestingContext) KlogConfigFolder() app.File {
	return app.NewFileOrPanic("/tmp/sample-klog-config-folder")
}

func (ctx *TestingContext) Meta() app.Meta {
	return app.Meta{
		Specification: "",
		License:       "",
		Version:       "v0.0",
		SrcHash:       "abc1234",
	}
}

func (ctx *TestingContext) ReadInputs(_ ...app.FileOrBookmarkName) ([]klog.Record, app.Error) {
	return ctx.records, nil
}

func (ctx *TestingContext) ReconcileFile(_ app.FileOrBookmarkName, creators []reconciling.Creator, reconcile ...reconciling.Reconcile) (*reconciling.Result, app.Error) {
	result, err := app.ApplyReconciler(ctx.records, ctx.blocks, creators, reconcile...)
	if err != nil {
		return nil, err
	}
	ctx.writtenFileContents = result.AllSerialised
	return result, nil
}

func (ctx *TestingContext) WriteFile(_ app.File, contents string) app.Error {
	ctx.writtenFileContents = contents
	return nil
}

func (ctx *TestingContext) Now() gotime.Time {
	return ctx.now
}

func (ctx *TestingContext) RetrieveTargetFile(fileArg app.FileOrBookmarkName) (app.FileWithContents, app.Error) {
	if fileArg == "" {
		return nil, app.NewError("Error", "Error", nil)
	}
	return app.NewFileWithContents(string(fileArg), "")
}

func (ctx *TestingContext) ReadBookmarks() (app.BookmarksCollection, app.Error) {
	return ctx.bookmarks, nil
}

func (ctx *TestingContext) ManipulateBookmarks(_ func(app.BookmarksCollection) app.Error) app.Error {
	return nil
}

func (ctx *TestingContext) Execute(cmd command.Command) app.Error {
	return ctx.execute(cmd)
}

func (ctx *TestingContext) Editors() (string, []command.Command) {
	return ctx.editorExplicit, ctx.editorsAuto
}

func (ctx *TestingContext) FileExplorers() []command.Command {
	return ctx.fileExplorers
}

func (ctx *TestingContext) Serialise() (tf.Styler, app.TextSerialiser) {
	return ctx.styler, ctx.serialiser
}

func (ctx *TestingContext) ConfigureSerialisation(fn func(tf.Styler, bool) (tf.Styler, bool)) {
	styler, decimalDuration := fn(ctx.styler, ctx.serialiser.DecimalDuration)
	ctx.styler = styler
	ctx.serialiser = app.NewSerialiser(styler, decimalDuration)
}

func (ctx *TestingContext) Debug(_ func()) {}

func (ctx *TestingContext) Config() app.Config {
	return *ctx.config
}
