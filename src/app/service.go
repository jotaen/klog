package app

import (
	"errors"
	"klog/record"
	"os"
	"os/exec"
)

type Service interface {
	Input() []record.Record
	Save([]record.Record) error
	CurrentFile() string
	BookmarkedFiles() []string
	OpenInEditor() error
	OpenInFileBrowser(string) error
	QuickStartAt(record.Date, record.Time) error
	QuickStopAt(record.Date, record.Time) error
}

type context struct {
	path *os.File
}

func NewService(destinationFilePath *os.File) Service {
	return &context{
		path: destinationFilePath,
	}
}

func (c *context) Input() []record.Record {
	return nil // TODO
}

func (c *context) Save([]record.Record) error {
	return nil // TODO
}

func (c *context) CurrentFile() string {
	return ""
}

func (c *context) BookmarkedFiles() []string {
	return nil
}

func (c *context) OpenInEditor() error {
	// open -t ...
	cmd := exec.Command("subl", c.path.Name())
	return cmd.Run()
}

func (c *context) OpenInFileBrowser(path string) error {
	cmd := exec.Command("open", path)
	return cmd.Run()
}

func (c *context) QuickStartAt(date record.Date, time record.Time) error {
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
	return c.Save(rs)
}

func (c *context) QuickStopAt(date record.Date, time record.Time) error {
	rs := c.Input()
	var recordToAlter *record.Record
	for _, r := range rs {
		if r.Date() == date && r.OpenRange() != nil {
			recordToAlter = &r
		}
	}
	if recordToAlter == nil {
		return errors.New("NO_OPEN_RANGE")
	}
	(*recordToAlter).StartOpenRange(time)
	return c.Save(rs)
}
