package app

import (
	"errors"
	"fmt"
	"io/ioutil"
	"klog"
	"klog/parser"
	"os"
	"os/exec"
	"os/user"
	"path/filepath"
	"strings"
)

var BinaryVersion string   // will be set during build
var BinaryBuildHash string // will be set during build

type Context interface {
	Print(string)
	HomeDir() string
	MetaInfo() struct {
		Version   string
		BuildHash string
	}
	RetrieveRecords(...string) ([]klog.Record, error)
	SetBookmark(string) error
	Bookmark() (File, error)
	OpenInFileBrowser(string) error
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

func (c *context) HomeDir() string {
	return c.homeDir
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
		if b, err := c.Bookmark(); err == nil {
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
		rs, errs := parser.Parse(content)
		if errs != nil {
			return nil, errs
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

func (c *context) Bookmark() (File, error) {
	bookmarkPath := c.HomeDir() + "/.klog/bookmark.klg"
	dest, err := os.Readlink(bookmarkPath)
	if err != nil {
		return File{}, err
	}
	return File{
		Name:     filepath.Base(dest),
		Location: filepath.Dir(dest),
		Path:     dest,
	}, nil
}

func (c *context) SetBookmark(path string) error {
	bookmark, err := filepath.Abs(path)
	if err != nil {
		return err
	}
	if !strings.HasSuffix(bookmark, ".klg") {
		return errors.New("File name must have .klg extension")
	}
	klogFolder := c.HomeDir() + "/.klog"
	err = os.MkdirAll(klogFolder, 0700)
	if err != nil {
		return err
	}
	symlink := klogFolder + "/bookmark.klg"
	_ = os.Remove(symlink)
	return os.Symlink(bookmark, symlink)
}

func (c *context) OpenInFileBrowser(path string) error {
	cmd := exec.Command("open", path)
	return cmd.Run()
}

func readFile(path string) (string, error) {
	contents, err := ioutil.ReadFile(path)
	if err != nil {
		return "", err
	}
	return string(contents), nil
}
