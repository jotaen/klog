package cli

import (
	"github.com/jotaen/klog/klog"
	"github.com/jotaen/klog/klog/app"
	"github.com/jotaen/klog/klog/app/cli/lib"
	"github.com/jotaen/klog/klog/app/cli/lib/command"
	"github.com/jotaen/klog/klog/app/cli/lib/terminalformat"
	"github.com/jotaen/klog/klog/parser"
	"github.com/jotaen/klog/klog/parser/reconciling"
	"github.com/jotaen/klog/klog/parser/txt"
	gotime "time"
)

func NewTestingContext() TestingContext {
	bc := app.NewEmptyBookmarksCollection()
	config := app.NewDefaultConfig()
	return TestingContext{
		State: State{
			printBuffer:         "",
			writtenFileContents: "",
		},
		now:            gotime.Now(),
		records:        nil,
		blocks:         nil,
		serialiser:     lib.CliSerialiser{},
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
	fileCfg := app.FromConfigFile{FileContents: configFile}
	err := fileCfg.Apply(ctx.config)
	if err != nil {
		panic(err)
	}
	return ctx
}

func (ctx TestingContext) _SetExecute(execute func(command.Command) app.Error) TestingContext {
	ctx.execute = execute
	return ctx
}

func (ctx TestingContext) _Run(cmd func(app.Context) app.Error) (State, app.Error) {
	cmdErr := cmd(&ctx)
	out := terminalformat.StripAllAnsiSequences(ctx.printBuffer)
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
	serialiser     parser.Serialiser
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

func (ctx *TestingContext) Serialiser() parser.Serialiser {
	return ctx.serialiser
}

func (ctx *TestingContext) SetSerialiser(s parser.Serialiser) {
	ctx.serialiser = s
}

func (ctx *TestingContext) Debug(_ func()) {}

func (ctx *TestingContext) Config() app.Config {
	return *ctx.config
}
