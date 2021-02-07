package app

import (
	"errors"
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
	SetBookmark(string) error
	Bookmark() *File
	OpenInFileBrowser(string) error
	AppendTemplateToFile(string, string) error
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
		if b := c.Bookmark(); b != nil {
			paths = []string{b.Path}
		} else {
			return nil, errors.New("No input file(s) specified; couldn’t read from bookmarked file either.")
		}
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

func (c *context) Bookmark() *File {
	bookmarkPath := c.KlogFolder() + "bookmark.klg"
	dest, err := os.Readlink(bookmarkPath)
	if err != nil {
		return nil
	}
	_, err = os.Stat(dest)
	if err != nil {
		return nil
	}
	return &File{
		Name:     filepath.Base(dest),
		Location: filepath.Dir(dest),
		Path:     dest,
	}
}

func (c *context) SetBookmark(path string) error {
	bookmark, err := filepath.Abs(path)
	if err != nil {
		return errors.New("Target file does not exist")
	}
	if !strings.HasSuffix(bookmark, ".klg") {
		return errors.New("File name must have .klg extension")
	}
	klogFolder := c.KlogFolder()
	err = os.MkdirAll(klogFolder, 0700)
	if err != nil {
		return errors.New("Unable to initialise ~/.klog folder")
	}
	symlink := klogFolder + "/bookmark.klg"
	_ = os.Remove(symlink)
	err = os.Symlink(bookmark, symlink)
	if err != nil {
		return errors.New("Failed to create bookmark")
	}
	return nil
}

func (c *context) OpenInFileBrowser(path string) error {
	cmd := exec.Command("open", path)
	return cmd.Run()
}

func (c *context) AppendTemplateToFile(filePath string, templateName string) error {
	location := c.KlogFolder() + templateName + ".template.klg"
	template, err := readFile(location)
	if err != nil {
		return errors.New("No such template: " + location)
	}
	instance, err := service.RenderTemplate(template, time.Now())
	if err != nil {
		return err
	}
	contents, err := readFile(filePath)
	if err != nil {
		return err
	}
	err = appendToFile(filePath, service.AppendableText(contents, instance))
	return err
}

func readFile(path string) (string, error) {
	contents, err := ioutil.ReadFile(path)
	if err != nil {
		return "", errors.New("Cannot read file: " + path)
	}
	return string(contents), nil
}

func appendToFile(path string, textToAppend string) error {
	file, err := os.OpenFile(path, os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		return errors.New("Cannot write to file: " + path)
	}
	defer file.Close()
	if _, err := file.WriteString(textToAppend); err != nil {
		return errors.New("Cannot write to file: " + path)
	}
	return nil
}
