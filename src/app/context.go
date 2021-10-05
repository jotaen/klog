package app

import (
	"bufio"
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

type FileOrBookmarkName string

type Context interface {
	Print(string)
	ReadLine() (string, Error)
	KlogFolder() string
	HomeFolder() string
	MetaInfo() struct {
		Version   string
		BuildHash string
	}
	ReadInputs(...FileOrBookmarkName) ([]Record, error)
	ReadFileInput(FileOrBookmarkName) (*parser.ParseResult, File, error)
	WriteFile(File, string) Error
	Now() gotime.Time
	ReadBookmarks() (BookmarksCollection, Error)
	ManipulateBookmarks(func(BookmarksCollection) Error) Error
	OpenInFileBrowser(File) Error
	OpenInEditor(FileOrBookmarkName, func(string)) Error
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

func (ctx *context) ReadLine() (string, Error) {
	scanner := bufio.NewScanner(os.Stdin)
	if scanner.Scan() {
		input := scanner.Text()
		return input, nil
	}
	return "", NewErrorWithCode(
		IO_ERROR,
		"Cannot process input",
		"Reading from stdin failed",
		nil,
	)
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

func (ctx *context) ReadInputs(fileArgs ...FileOrBookmarkName) ([]Record, error) {
	bc, bErr := ctx.ReadBookmarks()
	if bErr != nil {
		return nil, bErr
	}
	files, rErr := retrieveFirst([]Retriever{
		(&stdinRetriever{ReadStdin}).Retrieve,
		(&fileRetriever{ReadFile, bc}).Retrieve,
	}, fileArgs...)
	if rErr != nil {
		return nil, rErr
	}
	if len(files) == 0 {
		return nil, NewErrorWithCode(
			NO_INPUT_ERROR,
			"No input given",
			"Please do one of the following:\n"+
				"    a) specify one or multiple file names or bookmark names\n"+
				"    b) pipe file contents via stdin\n"+
				"    c) set a default bookmark to read from",
			nil,
		)
	}
	var records []Record
	for _, f := range files {
		pr, parserErrors := parser.Parse(f.content)
		if parserErrors != nil {
			return nil, parserErrors
		}
		records = append(records, pr.Records...)
	}
	return records, nil
}

func (ctx *context) retrieveTargetFile(fileArg FileOrBookmarkName) (*fileWithContent, Error) {
	bc, err := ctx.ReadBookmarks()
	if err != nil {
		return nil, err
	}
	inputs, err := (&fileRetriever{ReadFile, bc}).Retrieve(fileArg)
	if err != nil {
		return nil, err
	}
	if len(inputs) == 0 {
		return nil, NewErrorWithCode(
			NO_TARGET_FILE,
			"No file specified",
			"Either specify a file name or bookmark name, or set a default bookmark",
			nil,
		)
	}
	return inputs[0], nil
}

func (ctx *context) ReadFileInput(fileArg FileOrBookmarkName) (*parser.ParseResult, File, error) {
	target, err := ctx.retrieveTargetFile(fileArg)
	if err != nil {
		return nil, nil, err
	}
	pr, parserErrors := parser.Parse(target.content)
	if parserErrors != nil {
		return nil, nil, parserErrors
	}
	return pr, target, nil
}

func (ctx *context) WriteFile(target File, contents string) Error {
	if target == nil {
		panic("No path specified")
	}
	return WriteToFile(target, contents)
}

func (ctx *context) Now() gotime.Time {
	return gotime.Now()
}

func (ctx *context) initialiseKlogFolder() Error {
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
	bookmarksDatabase, err := ReadFile(ctx.bookmarkDatabasePath())
	if err != nil {
		// If database doesn’t exist, try to convert from legacy bookmark file.
		if os.IsNotExist(err.Original()) {
			legacyTarget, err := os.Readlink(ctx.bookmarkLegacySymlinkPath())
			if err == nil {
				bookmarksDatabase = `[{"name":"` + BOOKMARK_DEFAULT_NAME + `", "path": "` + legacyTarget + `"}]`
			}
		} else {
			return nil, err
		}
	}
	return NewBookmarksCollectionFromJson(bookmarksDatabase)
}

func (ctx *context) ManipulateBookmarks(manipulate func(BookmarksCollection) Error) Error {
	bc, bErr := ctx.ReadBookmarks()
	if bErr != nil {
		return bErr
	}
	mErr := manipulate(bc)
	if mErr != nil {
		return mErr
	}
	iErr := ctx.initialiseKlogFolder()
	if iErr != nil {
		return iErr
	}
	_ = os.Remove(ctx.bookmarkLegacySymlinkPath()) // Clean up legacy bookmark file, if exists
	return WriteToFile(ctx.bookmarkDatabasePath(), bc.ToJson())
}

func (ctx *context) bookmarkLegacySymlinkPath() string {
	return ctx.KlogFolder() + "bookmark.klg"
}

func (ctx *context) bookmarkDatabasePath() File {
	return NewFileOrPanic(ctx.KlogFolder() + "/bookmarks.json")
}

func (ctx *context) OpenInFileBrowser(target File) Error {
	cmd := exec.Command("open", target.Location())
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

func (ctx *context) OpenInEditor(fileArg FileOrBookmarkName, printHint func(string)) Error {
	target, err := ctx.retrieveTargetFile(fileArg)
	if err != nil {
		return err
	}
	hint := "You can specify your preferred editor via the $EDITOR environment variable.\n"
	preferredEditor := os.Getenv("EDITOR")
	editors := append([]string{preferredEditor}, POTENTIAL_EDITORS...)
	for _, editor := range editors {
		cmd := exec.Command(editor, target.Path())
		cmd.Stdin = os.Stdin
		cmd.Stdout = os.Stdout
		err := cmd.Run()
		if err == nil {
			if preferredEditor == "" {
				// Inform the user that they can configure their editor:
				printHint(hint)
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
	location := NewFileOrPanic(ctx.KlogFolder() + templateName + ".template.klg")
	template, err := ReadFile(location)
	if err != nil {
		return nil, NewError(
			"No such template",
			"There is no template at location "+location.Path(),
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
