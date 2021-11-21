/*
Package app contains the functionality that is related to the application layer.
This includes all code for the command line interface and the procedures to
interact with the runtime environment.
*/
package app

import (
	"bufio"
	"fmt"
	. "github.com/jotaen/klog/src"
	"github.com/jotaen/klog/src/parser"
	"github.com/jotaen/klog/src/parser/reconciling"
	"os"
	"os/exec"
	"os/user"
	"strings"
	gotime "time"
)

// FileOrBookmarkName is either a file name or a bookmark name
// as specified as argument on the command line.
type FileOrBookmarkName string

// Context is a representation of the runtime environment of klog.
// The commands carry out all side effects via this interface.
type Context interface {
	// Print prints to stdout.
	Print(string)

	// ReadLine reads user input from stdin.
	ReadLine() (string, Error)

	// KlogFolder returns the path of the .klog folder.
	KlogFolder() string

	// HomeFolder returns the path of the user’s home folder.
	HomeFolder() string

	// Meta returns miscellaneous meta information.
	Meta() Meta

	// ReadInputs retrieves all input from the given file or bookmark names.
	ReadInputs(...FileOrBookmarkName) ([]Record, error)

	// ReconcileFile ... // TODO
	ReconcileFile(FileOrBookmarkName, ...reconciling.Handler) error

	// Now returns the current timestamp.
	Now() gotime.Time

	// ReadBookmarks returns all configured bookmarks of the user.
	ReadBookmarks() (BookmarksCollection, Error)

	// ManipulateBookmarks saves a modified bookmark collection.
	ManipulateBookmarks(func(BookmarksCollection) Error) Error

	// OpenInFileBrowser tries to open the file explorer at the location of the file.
	OpenInFileBrowser(File) Error

	// OpenInEditor tries to open a file or bookmark in the user’s preferred $EDITOR.
	OpenInEditor(FileOrBookmarkName, func(string)) Error

	// Serialiser returns the current serialiser.
	Serialiser() *parser.Serialiser

	// SetSerialiser sets the current serialiser.
	SetSerialiser(*parser.Serialiser)
}

// Meta holds miscellaneous information about the klog binary.
type Meta struct {

	// Specification contains the file format specification in full text.
	Specification string

	// License contains the license text.
	License string

	// Version contains the build version.
	Version string

	// BuildHash contains a unique build identifier.
	BuildHash string
}

// NewContext creates a new Context object.
func NewContext(homeDir string, meta Meta, serialiser *parser.Serialiser) Context {
	if meta.Version == "" {
		meta.Version = "v?.?"
	}
	if meta.BuildHash == "" {
		strings.Repeat("?", 7)
	}
	return &context{
		homeDir,
		serialiser,
		meta,
	}
}

// NewContextFromEnv creates a Context object by automatically discovering certain parameters.
// It returns an error if the auto-discovery failed.
func NewContextFromEnv(meta Meta, serialiser *parser.Serialiser) (Context, error) {
	homeDir, err := user.Current()
	if err != nil {
		return nil, err
	}
	return NewContext(homeDir.HomeDir, meta, serialiser), nil
}

type context struct {
	homeDir    string
	serialiser *parser.Serialiser
	meta       Meta
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

func (ctx *context) Meta() Meta {
	return ctx.meta
}

func (ctx *context) ReadInputs(fileArgs ...FileOrBookmarkName) ([]Record, error) {
	bc, bErr := ctx.ReadBookmarks()
	if bErr != nil {
		return nil, bErr
	}
	files, rErr := retrieveFirst([]Retriever{
		(&StdinRetriever{ReadStdin}).Retrieve,
		(&FileRetriever{ReadFile, bc}).Retrieve,
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
	var allRecords []Record
	for _, f := range files {
		records, _, parserErrors := parser.Parse(f.Contents())
		if parserErrors != nil {
			return nil, parserErrors
		}
		allRecords = append(allRecords, records...)
	}
	return allRecords, nil
}

func (ctx *context) retrieveTargetFile(fileArg FileOrBookmarkName) (FileWithContents, Error) {
	bc, err := ctx.ReadBookmarks()
	if err != nil {
		return nil, err
	}
	inputs, err := (&FileRetriever{ReadFile, bc}).Retrieve(fileArg)
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

func (ctx *context) ReconcileFile(fileArg FileOrBookmarkName, handler ...reconciling.Handler) error {
	target, err := ctx.retrieveTargetFile(fileArg)
	if err != nil {
		return err
	}
	records, blocks, parserErrors := parser.Parse(target.Contents())
	if parserErrors != nil {
		return parserErrors
	}
	baseReconciler := reconciling.NewReconciler(records, blocks)
	result, rErr := reconciling.Chain(baseReconciler, handler...)
	if rErr != nil {
		return rErr
	}
	wErr := WriteToFile(target, result.FileContents())
	if wErr != nil {
		return wErr
	}
	ctx.Print("\n" + ctx.Serialiser().SerialiseRecords(result.Record()) + "\n")
	return nil
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
		// If that fails for whatever reason, don’t bother and create fresh collection.
		if os.IsNotExist(err.Original()) {
			newBc := NewEmptyBookmarksCollection()
			legacyTargetPath, rErr := os.Readlink(ctx.bookmarkLegacySymlinkPath())
			if rErr != nil {
				return newBc, nil
			}
			legacyTarget, fErr := NewFile(legacyTargetPath)
			if fErr != nil {
				return newBc, nil
			}
			newBc.Set(NewDefaultBookmark(legacyTarget))
			return newBc, nil
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
	return NewFileOrPanic(ctx.KlogFolder() + "bookmarks.json")
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

func (ctx *context) Serialiser() *parser.Serialiser {
	return ctx.serialiser
}

func (ctx *context) SetSerialiser(serialiser *parser.Serialiser) {
	if serialiser == nil {
		panic("Serialiser cannot be nil")
	}
	ctx.serialiser = serialiser
}
