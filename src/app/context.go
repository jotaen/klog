package app

import (
	"klog/record"
	"os/exec"
	"strings"
)

type Context interface {
	Config() Config
	BookmarkedFile() []record.Record
	OutputFilePath() string
	LatestFiles() []string
	OpenInEditor() error
	OpenInFileBrowser(string) error
}

type context struct {
	outputFilePath string
	config         Config
	history        []string
}

func NewContext(config Config, history []string) (Context, error) {
	return &context{
		config:  config,
		history: history,
	}, nil
}

func NewContextFromEnv() (Context, error) {
	config, err := func() (Config, error) {
		configToml, err := readFile("~/.klog.toml")
		if err != nil {
			return NewDefaultConfig(), nil
		}
		return NewConfigFromToml(configToml)
	}()
	if err != nil {
		return nil, err
	}

	history := func() []string {
		h, _ := readFile("~/.klog/history")
		hs := strings.Split(h, "\n")
		var result []string
		for _, x := range hs {
			result = append(result, strings.TrimSpace(x))
		}
		return result
	}()
	return NewContext(config, history)
}

func (c *context) Config() Config {
	return Config{} // TODO
}

func (c *context) BookmarkedFile() []record.Record {
	return nil // TODO
}

func (c *context) OutputFilePath() string {
	return c.outputFilePath
}

func (c *context) LatestFiles() []string {
	return c.history
}

func (c *context) OpenInEditor() error {
	// open -t ...
	cmd := exec.Command("subl", c.outputFilePath)
	return cmd.Run()
}

func (c *context) OpenInFileBrowser(path string) error {
	cmd := exec.Command("open", path)
	return cmd.Run()
}
