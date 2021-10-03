package app

import (
	"fmt"
	. "github.com/jotaen/klog/src"
	"github.com/jotaen/klog/src/parser"
	"github.com/jotaen/klog/src/parser/parsing"
	"os"
	"os/exec"
	"os/user"
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
	ReadBookmarks() (BookmarksCollection, Error)
	SaveBookmarks(BookmarksCollection) Error
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
	defaultBookmark func() (Bookmark, Error),
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
	b, err := defaultBookmark()
	if err != nil {
		return nil, err
	} else if b != nil {
		content, err := ReadFile(b.Target().Path)
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
	inputs, err := retrieveInputs(paths, ReadStdin, func() (Bookmark, Error) {
		bc, err := ctx.ReadBookmarks()
		if err != nil {
			return nil, err
		}
		return bc.Default(), nil
	})
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
		bc, err := ctx.ReadBookmarks()
		if err != nil {
			return nil, nil, err
		} else if bc.Default() == nil {
			return nil, nil, NewErrorWithCode(
				NO_TARGET_FILE,
				"No file specified",
				"You can either specify a file path, or you set a bookmark",
				nil,
			)
		}
		path = bc.Default().Target().Path
	}
	content, err := ReadFile(path)
	if err != nil {
		return nil, nil, err
	}
	pr, parserErrors := parser.Parse(content)
	if parserErrors != nil {
		return nil, nil, parserErrors
	}
	return pr, NewFile(path), nil
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

func (ctx *context) initializeKlogFolder() Error {
	klogFolder := ctx.KlogFolder()
	err := os.MkdirAll(klogFolder, 0700)
	flagAsHidden(klogFolder)
	if err != nil {
		return NewError(
			"Unable to initialise ~/.klog folder",
			"Please create a ~/.klog folder manually",
			err,
		)
	}
	return nil
}

func (ctx *context) ReadBookmarks() (BookmarksCollection, Error) {
	bookmarksDatabase, err := ReadFile(ctx.bookmarkDatabasePath().Path)
	if err != nil {
		// If database doesn’t exist, try to convert from legacy bookmark file.
		if os.IsNotExist(err.Original()) {
			legacyTarget, err := os.Readlink(ctx.bookmarkLegacySymlinkPath())
			if err == nil {
				bookmarksDatabase = `[{"name":"` + defaultName + `", "path": "` + legacyTarget + `"}]`
			}
		} else {
			return nil, err
		}
	}
	return NewBookmarksCollectionFromJson(bookmarksDatabase)
}

func (ctx *context) SaveBookmarks(bc BookmarksCollection) Error {
	err := ctx.initializeKlogFolder()
	if err != nil {
		return err
	}
	_ = os.Remove(ctx.bookmarkLegacySymlinkPath()) // Clean up legacy bookmark file, if exists
	return WriteToFile(ctx.bookmarkDatabasePath().Path, bc.ToJson())
}

func (ctx *context) bookmarkLegacySymlinkPath() string {
	return ctx.KlogFolder() + "bookmark.klg"
}

func (ctx *context) bookmarkDatabasePath() *File {
	return &File{
		"bookmarks.json",
		ctx.KlogFolder(),
		ctx.KlogFolder() + "/bookmarks.json",
	}
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
	hint := "You can specify your preferred editor via the $EDITOR environment variable.\n"
	preferredEditor := os.Getenv("EDITOR")
	editors := append([]string{preferredEditor}, POTENTIAL_EDITORS...)
	for _, editor := range editors {
		cmd := exec.Command(editor, path)
		cmd.Stdin = os.Stdin
		cmd.Stdout = os.Stdout
		err := cmd.Run()
		if err == nil {
			if preferredEditor == "" {
				// Inform the user that they can configure their editor:
				ctx.Print(hint)
			}
			return nil
		}
	}
	return NewError(
		"Cannot open any editor",
		hint,
		nil,
	)
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
