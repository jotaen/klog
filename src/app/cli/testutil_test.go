package cli

import (
	"github.com/jotaen/klog/lib/jotaen/terminalformat"
	. "github.com/jotaen/klog/src"
	"github.com/jotaen/klog/src/app"
	"github.com/jotaen/klog/src/app/cli/lib"
	"github.com/jotaen/klog/src/parser"
	"github.com/jotaen/klog/src/parser/parsing"
	gotime "time"
)

func NewTestingContext() TestingContext {
	return TestingContext{
		State: State{
			printBuffer:         "",
			writtenFileContents: "",
		},
		now:         gotime.Now(),
		records:     nil,
		parseResult: nil,
		serialiser:  lib.NewCliSerialiser(),
	}
}

func (ctx TestingContext) _SetRecords(records string) TestingContext {
	pr, err := parser.Parse(records)
	if err != nil {
		panic("Invalid records")
	}
	ctx.parseResult = pr
	ctx.records = pr.Records
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
	now         gotime.Time
	records     []Record
	parseResult *parser.ParseResult
	serialiser  *parser.Serialiser
}

func (ctx *TestingContext) Print(s string) {
	ctx.printBuffer += s
}

func (ctx *TestingContext) HomeFolder() string {
	return "~"
}

func (ctx *TestingContext) KlogFolder() string {
	return ctx.HomeFolder() + "/.klog/"
}

func (ctx *TestingContext) MetaInfo() struct {
	Version   string
	BuildHash string
} {
	return struct {
		Version   string
		BuildHash string
	}{"v0.0", "abcdef1"}
}

func (ctx *TestingContext) ReadInputs(_ ...string) ([]Record, error) {
	return ctx.records, nil
}

func (ctx *TestingContext) ReadFileInput(string) (*parser.ParseResult, *app.File, error) {
	return ctx.parseResult, nil, nil
}

func (ctx *TestingContext) WriteFile(_ *app.File, contents string) app.Error {
	ctx.writtenFileContents = contents
	return nil
}

func (ctx *TestingContext) Now() gotime.Time {
	return ctx.now
}

func (ctx *TestingContext) SetBookmark(_ string) app.Error {
	return nil
}

func (ctx *TestingContext) Bookmark() (*app.File, app.Error) {
	return &app.File{
		Name:     "myfile.klg",
		Location: "/",
		Path:     "/myfile.klg",
	}, nil
}

func (ctx *TestingContext) UnsetBookmark() app.Error {
	return nil
}

func (ctx *TestingContext) OpenInFileBrowser(_ string) app.Error {
	return nil
}

func (ctx *TestingContext) OpenInEditor(_ string) app.Error {
	return nil
}

func (ctx *TestingContext) InstantiateTemplate(_ string) ([]parsing.Text, app.Error) {
	return nil, nil
}

func (ctx *TestingContext) Serialiser() *parser.Serialiser {
	return ctx.serialiser
}

func (ctx *TestingContext) SetSerialiser(serialiser *parser.Serialiser) {
	if serialiser == nil {
		panic("Serialiser cannot be nil")
	}
	ctx.serialiser = serialiser
}
