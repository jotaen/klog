package app

import (
	"fmt"
	. "klog"
	"klog/parser"
	"klog/parser/parsing"
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
	ReadInputs(...string) ([]Record, error)
	ReadFileInput(string) (*parser.ParseResult, *File, error)
	WriteFile(*File, string) Error
	Now() gotime.Time
	Bookmark() (*File, Error)
	SetBookmark(string) Error
	UnsetBookmark() Error
	OpenInFileBrowser(string) Error
	OpenInEditor(string) Error
	InstantiateTemplate(string) ([]parsing.Text, Error)
	Serialiser() *parser.Serialiser
	SetSerialiser(*parser.Serialiser)
}

type context struct {
	homeDir    string
	serialiser *parser.Serialiser
}

func NewContext(homeDir string, serialiser *parser.Serialiser) (Context, error) {
	return &context{
		homeDir:    homeDir,
		serialiser: serialiser,
	}, nil
}

func NewContextFromEnv(serialiser *parser.Serialiser) (Context, error) {
	homeDir, err := user.Current()
	if err != nil {
		return nil, err
	}
	return NewContext(homeDir.HomeDir, serialiser)
}

func (ctx *context) Print(text string) {
	fmt.Print(text)
}

func (ctx *context) HomeFolder() string {
	return ctx.homeDir
}

func (ctx *context) KlogFolder() string {
	return ctx.homeDir + "/.klog/"
}

func (ctx *context) MetaInfo() struct {
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
			content, err := ReadFile(p)
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
		content, err := ReadFile(b.Path)
		if err != nil {
			return nil, err
		}
		return []string{content}, nil
	}
	return nil, NewErrorWithCode(
		NO_INPUT_ERROR,
		"No input given",
		"Please do one of the following:\n"+
			"    a) pass one or multiple file names as argument\n"+
			"    b) pipe file contents via stdin\n"+
			"    c) specify a bookmark to read from by default",
		err,
	)
}

func (ctx *context) ReadInputs(paths ...string) ([]Record, error) {
	inputs, err := retrieveInputs(paths, ReadStdin, ctx.bookmarkOrNil)
	if err != nil {
		return nil, err
	}
	var records []Record
	for _, in := range inputs {
		pr, parserErrors := parser.Parse(in)
		if parserErrors != nil {
			return nil, parserErrors
		}
		records = append(records, pr.Records...)
	}
	return records, nil
}

func (ctx *context) ReadFileInput(path string) (*parser.ParseResult, *File, error) {
	if path == "" {
		b, err := ctx.bookmarkOrNil()
		if err != nil {
			return nil, nil, err
		} else if b == nil {
			return nil, nil, NewErrorWithCode(
				NO_TARGET_FILE,
				"No file specified",
				"You can either specify a file path, or you set a bookmark",
				nil,
			)
		}
		path = b.Path
	}
	content, err := ReadFile(path)
	if err != nil {
		return nil, nil, err
	}
	pr, parserErrors := parser.Parse(content)
	if parserErrors != nil {
		return nil, nil, parserErrors
	}
	return pr, newFile(path), nil
}

func (ctx *context) WriteFile(target *File, contents string) Error {
	if target == nil {
		panic("No path specified")
	}
	return WriteToFile(target.Path, contents)
}

func (ctx *context) Now() gotime.Time {
	return gotime.Now()
}

type File struct {
	Name     string
	Location string
	Path     string
}

func newFile(path string) *File {
	return &File{
		Name:     filepath.Base(path),
		Location: filepath.Dir(path),
		Path:     path,
	}
}

func (ctx *context) bookmarkOrNil() (*File, Error) {
	bookmarkPath := ctx.bookmarkOrigin()
	dest, err := os.Readlink(bookmarkPath)
	if err != nil {
		return nil, nil
	}
	_, err = os.Stat(dest)
	if err != nil {
		return nil, NewErrorWithCode(
			BOOKMARK_ACCESS_ERROR,
			"Bookmark doesn’t point to valid file",
			"Please check the current bookmark location or set a new one",
			err,
		)
	}
	return newFile(dest), nil
}

func (ctx *context) bookmarkOrigin() string {
	return ctx.KlogFolder() + "bookmark.klg"
}

func (ctx *context) Bookmark() (*File, Error) {
	b, err := ctx.bookmarkOrNil()
	if err != nil {
		return nil, err
	}
	if b == nil {
		return nil, NewErrorWithCode(
			NO_BOOKMARK_SET_ERROR,
			"No bookmark set",
			"You can set a bookmark by running: klog bookmark set somefile.klg",
			err,
		)
	}
	return b, nil
}

func (ctx *context) SetBookmark(path string) Error {
	bookmark, err := filepath.Abs(path)
	if err != nil {
		return NewErrorWithCode(
			BOOKMARK_ACCESS_ERROR,
			"Invalid target file",
			"Please check the file path",
			err,
		)
	}
	klogFolder := ctx.KlogFolder()
	err = os.MkdirAll(klogFolder, 0700)
	flagAsHidden(klogFolder)
	if err != nil {
		return NewError(
			"Unable to initialise ~/.klog folder",
			"Please create a ~/.klog folder manually",
			err,
		)
	}
	appErr := ctx.UnsetBookmark()
	if appErr != nil {
		return appErr
	}
	appErr = createSymlinkForBookmark(bookmark, ctx.bookmarkOrigin())
	if appErr != nil {
		return appErr
	}
	return nil
}

func (ctx *context) UnsetBookmark() Error {
	err := os.Remove(ctx.bookmarkOrigin())
	if err != nil && !os.IsNotExist(err) {
		return NewErrorWithCode(
			BOOKMARK_ACCESS_ERROR,
			"Failed to unset bookmark",
			"The current bookmark could not be cleared",
			err,
		)
	}
	return nil
}

func (ctx *context) OpenInFileBrowser(path string) Error {
	cmd := exec.Command("open", path)
	err := cmd.Run()
	if err != nil {
		return NewError(
			"Failed to open file browser",
			err.Error(),
			err,
		)
	}
	return nil
}

func (ctx *context) OpenInEditor(path string) Error {
	editor := os.Getenv("EDITOR")
	if editor == "" {
		return NewError(
			"No default editor set",
			"Please specify you editor via the $EDITOR environment variable",
			nil,
		)
	}
	cmd := exec.Command(editor, path)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	err := cmd.Run()
	if err != nil {
		return NewError(
			"Cannot open editor",
			"Tried to run: "+editor+" "+path,
			err,
		)
	}
	return nil
}

func (ctx *context) InstantiateTemplate(templateName string) ([]parsing.Text, Error) {
	location := ctx.KlogFolder() + templateName + ".template.klg"
	template, err := ReadFile(location)
	if err != nil {
		return nil, NewError(
			"No such template",
			"There is no template at location "+location,
			err,
		)
	}
	instance, tErr := parser.RenderTemplate(template, ctx.Now())
	if tErr != nil {
		return nil, NewError(
			"Invalid template",
			tErr.Error(),
			tErr,
		)
	}
	return instance, nil
}

func (ctx *context) Serialiser() *parser.Serialiser {
	return ctx.serialiser
}

func (ctx *context) SetSerialiser(serialiser *parser.Serialiser) {
	if serialiser == nil {
		panic("Serialiser cannot be nil")
	}
	ctx.serialiser = serialiser
}
