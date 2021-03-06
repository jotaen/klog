package app

import (
	"fmt"
	"klog"
	"klog/parser"
	"klog/service"
	"os"
	"os/exec"
	"os/user"
	"path/filepath"
	"strings"
	gotime "time"
)

var BinaryVersion string   // will be set during build
var BinaryBuildHash string // will be set during build

type Context interface {
	Print(string)
	KlogFolder() string
	HomeFolder() string
	MetaInfo() struct {
		Version   string
		BuildHash string
	}
	RetrieveRecords(...string) ([]klog.Record, error)
	Now() gotime.Time
	Bookmark() (*File, Error)
	SetBookmark(string) Error
	UnsetBookmark() Error
	OpenInFileBrowser(string) Error
	OpenInEditor(string) Error
	AppendTemplateToFile(string, string) Error
}

type context struct {
	homeDir string
}

func NewContext(homeDir string) (Context, error) {
	return &context{
		homeDir: homeDir,
	}, nil
}

func NewContextFromEnv() (Context, error) {
	homeDir, err := user.Current()
	if err != nil {
		return nil, err
	}
	return NewContext(homeDir.HomeDir)
}

func (c *context) Print(text string) {
	fmt.Print(text)
}

func (c *context) HomeFolder() string {
	return c.homeDir
}

func (c *context) KlogFolder() string {
	return c.homeDir + "/.klog/"
}

func (c *context) MetaInfo() struct {
	Version   string
	BuildHash string
} {
	return struct {
		Version   string
		BuildHash string
	}{
		Version: func() string {
			if BinaryVersion == "" {
				return "v?.?"
			}
			return BinaryVersion
		}(),
		BuildHash: func() string {
			if BinaryBuildHash == "" {
				return strings.Repeat("?", 7)
			}
			if len(BinaryBuildHash) > 7 {
				return BinaryBuildHash[:7]
			}
			return BinaryBuildHash
		}(),
	}
}

func retrieveInputs(
	filePaths []string,
	readStdin func() (string, Error),
	bookmarkOrNil func() (*File, Error),
) ([]string, Error) {
	if len(filePaths) > 0 {
		var result []string
		for _, p := range filePaths {
			content, err := readFile(p)
			if err != nil {
				return nil, err
			}
			result = append(result, content)
		}
		return result, nil
	}
	stdin, err := readStdin()
	if err != nil {
		return nil, err
	}
	if stdin != "" {
		return []string{stdin}, nil
	}
	b, err := bookmarkOrNil()
	if err != nil {
		return nil, err
	} else if b != nil {
		content, err := readFile(b.Path)
		if err != nil {
			return nil, err
		}
		return []string{content}, nil
	}
	return nil, appError{
		"No input given",
		"Please do one of the following:\n" +
			"    a) pass one or multiple file names as argument\n" +
			"    b) pipe file contents via stdin\n" +
			"    c) specify a bookmark to read from by default",
	}
}

func (c *context) RetrieveRecords(paths ...string) ([]klog.Record, error) {
	inputs, err := retrieveInputs(paths, readStdin, c.bookmarkOrNil)
	if err != nil {
		return nil, err
	}
	var records []klog.Record
	for _, in := range inputs {
		pr, parserErrors := parser.Parse(in)
		if parserErrors != nil {
			return nil, parserErrors
		}
		records = append(records, pr.Records...)
	}
	return records, nil
}

func (c *context) Now() gotime.Time {
	return gotime.Now()
}

type File struct {
	Name     string
	Location string
	Path     string
}

func (c *context) bookmarkOrNil() (*File, Error) {
	bookmarkPath := c.bookmarkOrigin()
	dest, err := os.Readlink(bookmarkPath)
	if err != nil {
		return nil, nil
	}
	_, err = os.Stat(dest)
	if err != nil {
		return nil, appError{
			"Bookmark doesn’t point to valid file",
			"Please check the current bookmark location or set a new one",
		}
	}
	return &File{
		Name:     filepath.Base(dest),
		Location: filepath.Dir(dest),
		Path:     dest,
	}, nil
}

func (c *context) bookmarkOrigin() string {
	return c.KlogFolder() + "bookmark.klg"
}

func (c *context) Bookmark() (*File, Error) {
	b, err := c.bookmarkOrNil()
	if err != nil {
		return nil, err
	}
	if b == nil {
		return nil, appError{
			"No bookmark set",
			"You can set a bookmark by running: klog bookmark set somefile.klg",
		}
	}
	return b, nil
}

func (c *context) SetBookmark(path string) Error {
	bookmark, err := filepath.Abs(path)
	if err != nil {
		return appError{
			"Invalid target file",
			"Please check the file path",
		}
	}
	klogFolder := c.KlogFolder()
	err = os.MkdirAll(klogFolder, 0700)
	if err != nil {
		return appError{
			"Unable to initialise ~/.klog folder",
			"Please create a ~/.klog folder manually",
		}
	}
	symlink := c.bookmarkOrigin()
	_ = os.Remove(symlink)
	err = os.Symlink(bookmark, symlink)
	if err != nil {
		return appError{
			"Failed to create bookmark",
			"",
		}
	}
	return nil
}

func (c *context) UnsetBookmark() Error {
	return removeFile(c.bookmarkOrigin())
}

func (c *context) OpenInFileBrowser(path string) Error {
	cmd := exec.Command("open", path)
	err := cmd.Run()
	if err != nil {
		return appError{
			"Failed to open file browser",
			err.Error(),
		}
	}
	return nil
}

func (c *context) OpenInEditor(path string) Error {
	editor := os.Getenv("EDITOR")
	if editor == "" {
		return appError{
			"No default editor set",
			"Please specify you editor via the $EDITOR environment variable",
		}
	}
	cmd := exec.Command(editor, path)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	err := cmd.Run()
	if err != nil {
		return appError{
			"Cannot open editor",
			"Tried to run: " + editor + " " + path,
		}
	}
	return nil
}

func (c *context) AppendTemplateToFile(filePath string, templateName string) Error {
	location := c.KlogFolder() + templateName + ".template.klg"
	template, err := readFile(location)
	if err != nil {
		return appError{
			"No such template",
			"There is no template at location " + location,
		}
	}
	instance, tErr := service.RenderTemplate(template, c.Now())
	if tErr != nil {
		return appError{
			"Invalid template",
			tErr.Error(),
		}
	}
	contents, err := readFile(filePath)
	if err != nil {
		return err
	}
	return appendToFile(filePath, service.AppendableText(contents, instance))
}
