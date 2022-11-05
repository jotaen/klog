package cli

import (
	"github.com/jotaen/klog/klog"
	"github.com/jotaen/klog/klog/app"
	"github.com/jotaen/klog/klog/app/cli/lib"
	"github.com/jotaen/klog/klog/app/cli/lib/terminalformat"
	"github.com/jotaen/klog/klog/parser"
	"github.com/jotaen/klog/klog/parser/reconciling"
	"github.com/jotaen/klog/klog/parser/txt"
	gotime "time"
)

func NewTestingContext() TestingContext {
	bc := app.NewEmptyBookmarksCollection()
	return TestingContext{
		State: State{
			printBuffer:         "",
			writtenFileContents: "",
		},
		now:        gotime.Now(),
		records:    nil,
		blocks:     nil,
		serialiser: lib.CliSerialiser{},
		bookmarks:  bc,
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

func (ctx TestingContext) _Run(cmd func(app.Context) error) (State, error) {
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
	now        gotime.Time
	records    []klog.Record
	blocks     []txt.Block
	serialiser parser.Serialiser
	bookmarks  app.BookmarksCollection
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

func (ctx *TestingContext) KlogFolder() string {
	return ctx.HomeFolder() + "/.klog/"
}

func (ctx *TestingContext) Meta() app.Meta {
	return app.Meta{
		Specification: "",
		License:       "",
		Version:       "v0.0",
		BuildHash:     "abcdef1",
	}
}

func (ctx *TestingContext) ReadInputs(_ ...app.FileOrBookmarkName) ([]klog.Record, app.Error) {
	return ctx.records, nil
}

func (ctx *TestingContext) ReconcileFile(doWrite bool, _ app.FileOrBookmarkName, creators []reconciling.Creator, reconcile reconciling.Reconcile) (*reconciling.Result, app.Error) {
	result, err := app.ApplyReconciler(ctx.records, ctx.blocks, creators, reconcile)
	if err != nil {
		return nil, err
	}
	if doWrite {
		ctx.writtenFileContents = result.AllSerialised
	}
	return result, nil
}

func (ctx *TestingContext) WriteFile(_ app.File, contents string) app.Error {
	ctx.writtenFileContents = contents
	return nil
}

func (ctx *TestingContext) Now() gotime.Time {
	return ctx.now
}

func (ctx *TestingContext) ReadBookmarks() (app.BookmarksCollection, app.Error) {
	return ctx.bookmarks, nil
}

func (ctx *TestingContext) ManipulateBookmarks(_ func(app.BookmarksCollection) app.Error) app.Error {
	return nil
}

func (ctx *TestingContext) OpenInFileBrowser(_ app.FileOrBookmarkName) app.Error {
	return nil
}

func (ctx *TestingContext) OpenInEditor(_ app.FileOrBookmarkName, _ func(string)) app.Error {
	return nil
}

func (ctx *TestingContext) Serialiser() parser.Serialiser {
	return ctx.serialiser
}

func (ctx *TestingContext) SetSerialiser(s parser.Serialiser) {
	ctx.serialiser = s
}

func (ctx *TestingContext) Debug(_ func()) {}

func (ctx *TestingContext) Preferences() app.Preferences {
	return app.NewDefaultPreferences()
}
