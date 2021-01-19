package app

import (
	"klog/parser"
	"klog/record"
	"os/exec"
	"os/user"
	"strings"
)

type Context struct {
	bookmarkedFile string
	config         Config
	history        []string
	homeDir        string
}

func NewContext(config Config, history []string) (*Context, error) {
	homeDir, err := user.Current()
	if err != nil {
		return nil, err
	}
	return &Context{
		config:  config,
		history: history,
		homeDir: homeDir.HomeDir,
	}, nil
}

func NewContextFromEnv() (*Context, error) {
	config, err := func() (Config, error) {
		configToml, err := ReadFile("~/.klog.toml")
		if err != nil {
			return NewDefaultConfig(), nil
		}
		return NewConfigFromToml(configToml)
	}()
	if err != nil {
		return nil, err
	}

	history := func() []string {
		h, _ := ReadFile("~/.klog/history")
		hs := strings.Split(h, "\n")
		var result []string
		for _, x := range hs {
			result = append(result, strings.TrimSpace(x))
		}
		return result
	}()
	return NewContext(config, history)
}

func (c *Context) HomeDir() string {
	return c.homeDir
}

func (c *Context) RetrieveRecords(paths []string) ([]record.Record, error) {
	var records []record.Record
	for _, p := range paths {
		content, err := ReadFile(p)
		if err != nil {
			return nil, err
		}
		rs, errs := parser.Parse(content)
		if errs != nil {
			return nil, errs
		}
		records = append(records, rs...)
	}
	return records, nil
}

func (c *Context) Config() Config {
	return Config{} // TODO
}

func (c *Context) BookmarkedFile() []record.Record {
	return nil // TODO
}

func (c *Context) LatestFiles() []string {
	return c.history
}

func (c *Context) OpenInFileBrowser(path string) error {
	cmd := exec.Command("open", path)
	return cmd.Run()
}
