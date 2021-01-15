package app

import (
	"errors"
	"klog/parser"
	"klog/parser/engine"
	"klog/record"
	"os/exec"
	"strings"
)

type Service interface {
	Input() []record.Record
	SetInput(string) error
	Save([]record.Record) error
	OutputFilePath() string
	LatestFiles() []string
	OpenInEditor() error
	OpenInFileBrowser(string) error
	QuickStartAt(record.Date, record.Time) (record.Record, error)
	QuickStopAt(record.Date, record.Time) (record.Record, error)
}

type context struct {
	input          []record.Record
	outputFilePath string
	config         Config
	history        []string
}

func NewService(configToml string, history []string) (Service, error) {
	cfg, err := NewConfigFromToml(configToml)
	if err != nil {
		return nil, err
	}
	return &context{
		config:  cfg,
		history: history,
	}, nil
}

func NewServiceWithConfigFiles() (Service, error) {
	configToml, err := readFile("~/.klog.toml")
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
	return NewService(configToml, history)
}

func (c *context) Input() []record.Record {
	return c.input
}

func (c *context) SetInput(recordsText string) error {
	rs, err := parser.Parse(recordsText)
	if err != nil {
		return err
	}
	c.input = rs
	return nil
}

func (c *context) Save(rs []record.Record) error {
	if c.outputFilePath == "" {
		return errors.New("NO_OUTPUT_TARGET")
	}
	return writeFile(c.outputFilePath, parser.ToPlainText(rs))
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

func (c *context) QuickStartAt(date record.Date, time record.Time) (record.Record, error) {
	rs := c.Input()
	var recordToAlter *record.Record
	for _, r := range rs {
		if r.Date() == date {
			recordToAlter = &r
		}
	}
	if recordToAlter == nil {
		r := record.NewRecord(date)
		recordToAlter = &r
	}
	(*recordToAlter).StartOpenRange(time)
	err := c.Save(rs)
	return *recordToAlter, err
}

func (c *context) QuickStopAt(date record.Date, time record.Time) (record.Record, error) {
	rs := c.Input()
	var recordToAlter *record.Record
	for _, r := range rs {
		if r.Date() == date && r.OpenRange() != nil {
			recordToAlter = &r
		}
	}
	if recordToAlter == nil {
		return nil, errors.New("NO_OPEN_RANGE")
	}
	(*recordToAlter).StartOpenRange(time)
	err := c.Save(rs)
	return *recordToAlter, err
}
