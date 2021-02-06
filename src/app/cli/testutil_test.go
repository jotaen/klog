package cli

import (
	"klog"
	"klog/app"
	"klog/parser"
	"regexp"
)

var ansiSequencePattern = regexp.MustCompile(`\x1b\[[\d;]+m`)

func RunWithContext(records string, cmd func(app.Context) error) (string, error) {
	rs, err := parser.Parse(records)
	if err != nil {
		panic("Invalid records")
	}
	ctx := &TestContext{
		records: rs,
	}
	cmdErr := cmd(ctx)
	out := ansiSequencePattern.ReplaceAllString(ctx.printBuffer, "")
	if len(out) > 0 && out[0] != '\n' {
		out = "\n" + out
	}
	return out, cmdErr
}

type TestContext struct {
	printBuffer string
	records     []klog.Record
}

func (m *TestContext) Print(s string) {
	m.printBuffer += s
}

func (m *TestContext) HomeDir() string {
	return "~"
}

func (m *TestContext) MetaInfo() struct {
	Version   string
	BuildHash string
} {
	return struct {
		Version   string
		BuildHash string
	}{"v0.0", "abcdef1"}
}

func (m *TestContext) RetrieveRecords(_ ...string) ([]klog.Record, error) {
	return m.records, nil
}

func (m *TestContext) SetBookmark(_ string) error {
	return nil
}

func (m *TestContext) Bookmark() (app.File, error) {
	return app.File{
		Name:     "myfile.klg",
		Location: "/",
		Path:     "/myfile.klg",
	}, nil
}

func (m *TestContext) OpenInFileBrowser(_ string) error {
	return nil
}
