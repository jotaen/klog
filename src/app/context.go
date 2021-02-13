package app

import (
	"fmt"
	"io/ioutil"
	"klog"
	"klog/parser"
	"klog/service"
	"os"
	"os/exec"
	"os/user"
	"path/filepath"
	"strings"
	"time"
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
	SetBookmark(string) Error
	Bookmark() (*File, Error)
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

func (c *context) RetrieveRecords(paths ...string) ([]klog.Record, error) {
	if len(paths) == 0 {
		b, err := c.Bookmark()
		if err != nil {
			return nil, appError{
				"No input files specified",
				"Either specify input files, or set a bookmark",
			}
		}
		paths = []string{b.Path}
	}
	var records []klog.Record
	for _, p := range paths {
		content, err := readFile(p)
		if err != nil {
			return nil, err
		}
		rs, parserErrors := parser.Parse(content)
		if parserErrors != nil {
			return nil, parserErrors
		}
		records = append(records, rs...)
	}
	return records, nil
}

type File struct {
	Name     string
	Location string
	Path     string
}

func (c *context) Bookmark() (*File, Error) {
	bookmarkPath := c.KlogFolder() + "bookmark.klg"
	dest, err := os.Readlink(bookmarkPath)
	if err != nil {
		return nil, appError{
			"No bookmark set",
			"You can set a bookmark by running: klog bookmark set somefile.klg",
		}
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
	symlink := klogFolder + "/bookmark.klg"
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
	instance, tErr := service.RenderTemplate(template, time.Now())
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

func readFile(path string) (string, Error) {
	contents, err := ioutil.ReadFile(path)
	if err != nil {
		return "", appError{
			"Cannot read file",
			"Location: " + path,
		}
	}
	return string(contents), nil
}

func appendToFile(path string, textToAppend string) Error {
	file, err := os.OpenFile(path, os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		return appError{
			"Cannot write to file",
			"Location: " + path,
		}
	}
	defer file.Close()
	if _, err := file.WriteString(textToAppend); err != nil {
		return appError{
			"Cannot write to file",
			"Location: " + path,
		}
	}
	return nil
}
