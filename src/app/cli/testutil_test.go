package cli

import (
	. "klog"
	"klog/app"
	"klog/parser"
	"klog/parser/parsing"
	"regexp"
	gotime "time"
)

var ansiSequencePattern = regexp.MustCompile(`\x1b\[[\d;]+m`)

func NewTestingContext() TestingContext {
	return TestingContext{
		State: State{
			printBuffer: "",
			writtenFileContents: "",
		},
		now:         gotime.Now(),
		records:     nil,
		parseResult: nil,
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
	out := ansiSequencePattern.ReplaceAllString(ctx.printBuffer, "")
	if len(out) > 0 && out[0] != '\n' {
		out = "\n" + out
	}
	return State{out, ctx.writtenFileContents}, cmdErr
}

type State struct {
	printBuffer string
	writtenFileContents string
}

type TestingContext struct {
	State
	now         gotime.Time
	records     []Record
	parseResult *parser.ParseResult
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

func (ctx *TestingContext) ReadFileInput(string) (*parser.ParseResult, error) {
	return ctx.parseResult, nil
}

func (ctx *TestingContext) WriteFile(_ string, contents string) app.Error {
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
