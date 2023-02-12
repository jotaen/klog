/*
Package app contains the functionality that is related to the application layer.
This includes all code for the command line interface and the procedures to
interact with the runtime environment.
*/
package app

import (
	"bufio"
	"fmt"
	"github.com/jotaen/klog/klog"
	"github.com/jotaen/klog/klog/app/cli/lib/command"
	"github.com/jotaen/klog/klog/parser"
	"github.com/jotaen/klog/klog/parser/reconciling"
	"github.com/jotaen/klog/klog/parser/txt"
	"os"
	"os/exec"
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
	ReadInputs(...FileOrBookmarkName) ([]klog.Record, Error)

	// RetrieveTargetFile returns the desired file, requiring that there is exactly one.
	RetrieveTargetFile(fileArg FileOrBookmarkName) (FileWithContents, Error)

	// ReconcileFile applies one or more reconcile handlers to a file and saves it.
	ReconcileFile(FileOrBookmarkName, []reconciling.Creator, reconciling.Reconcile) (*reconciling.Result, Error)

	// Now returns the current timestamp.
	Now() gotime.Time

	// ReadBookmarks returns all configured bookmarks of the user.
	ReadBookmarks() (BookmarksCollection, Error)

	// ManipulateBookmarks saves a modified bookmark collection.
	ManipulateBookmarks(func(BookmarksCollection) Error) Error

	// Execute attempts to run a command on the system.
	Execute(command.Command) Error

	// Editors returns commands to launch a text editor on the system.
	// - The string is a user-specified command, if specified.
	// - The command list is a prioritised list of predefined commands.
	Editors() (string, []command.Command)

	// FileExplorers returns commands to launch a file explorer on the system.
	FileExplorers() []command.Command

	// Serialiser returns the current serialiser.
	Serialiser() parser.Serialiser

	// SetSerialiser sets a new serialiser.
	SetSerialiser(parser.Serialiser)

	// Debug takes a void function that is only executed in debug mode.
	Debug(func())

	// Preferences returns the current preferences.
	Preferences() Preferences
}

// Meta holds miscellaneous information about the klog binary.
type Meta struct {

	// Specification contains the file format specification in full text.
	Specification string

	// License contains the license text.
	License string

	// Version contains the release version, e.g. `v2.7`.
	Version string

	// SrcHash contains the hash of the sources that the binary was built from.
	SrcHash string
}

// Preferences are user-defined configuration.
type Preferences struct {
	IsDebug    bool
	Editor     string
	NoColour   bool
	CpuKernels int // Must be 1 or higher
}

func NewDefaultPreferences() Preferences {
	return Preferences{
		IsDebug:    false,
		Editor:     "",
		NoColour:   false,
		CpuKernels: 1,
	}
}

// NewContext creates a new Context object.
func NewContext(homeDir string, meta Meta, serialiser parser.Serialiser, prefs Preferences) Context {
	parserEngine := parser.NewSerialParser()
	if prefs.CpuKernels > 1 {
		parserEngine = parser.NewParallelParser(prefs.CpuKernels)
	}
	return &context{
		homeDir,
		parserEngine,
		serialiser,
		meta,
		prefs,
	}
}

type context struct {
	homeDir    string
	parser     parser.Parser
	serialiser parser.Serialiser
	meta       Meta
	prefs      Preferences
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

func (ctx *context) ReadInputs(fileArgs ...FileOrBookmarkName) ([]klog.Record, Error) {
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
	var allRecords []klog.Record
	for _, f := range files {
		records, _, errs := ctx.parser.Parse(f.Contents())
		if errs != nil {
			return nil, NewParserErrors(errs)
		}
		allRecords = append(allRecords, records...)
	}
	return allRecords, nil
}

func (ctx *context) RetrieveTargetFile(fileArg FileOrBookmarkName) (FileWithContents, Error) {
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

func (ctx *context) ReconcileFile(fileArg FileOrBookmarkName, creators []reconciling.Creator, reconcile reconciling.Reconcile) (*reconciling.Result, Error) {
	target, err := ctx.RetrieveTargetFile(fileArg)
	if err != nil {
		return nil, err
	}
	records, blocks, errs := ctx.parser.Parse(target.Contents())
	if errs != nil {
		return nil, NewParserErrors(errs)
	}
	result, aErr := ApplyReconciler(records, blocks, creators, reconcile)
	if aErr != nil {
		return nil, aErr
	}
	wErr := WriteToFile(target, result.AllSerialised)
	if wErr != nil {
		return nil, wErr
	}
	return result, nil
}

func ApplyReconciler(records []klog.Record, blocks []txt.Block, creators []reconciling.Creator, reconcile reconciling.Reconcile) (*reconciling.Result, Error) {
	reconciler := func() *reconciling.Reconciler {
		for _, createReconciler := range creators {
			// Both the creator and the created reconciler might be nil,
			// to indicate it’s not eligible.
			if createReconciler == nil {
				continue
			}
			r := createReconciler(records, blocks)
			if r != nil {
				return r
			}
		}
		return nil
	}()
	if reconciler == nil {
		return nil, NewErrorWithCode(
			LOGICAL_ERROR,
			"No such record",
			"Please create or specify a record for this operation",
			nil,
		)
	}
	result, rErr := reconcile(reconciler)
	if rErr != nil {
		return nil, NewErrorWithCode(
			LOGICAL_ERROR,
			"Manipulation failed",
			rErr.Error(),
			rErr,
		)
	}
	return result, nil
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
		if os.IsNotExist(err.Original()) {
			// An absent bookmarks file is equivalent to an empty one.
			return NewEmptyBookmarksCollection(), nil
		}
		return nil, err
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
	return WriteToFile(ctx.bookmarkDatabasePath(), bc.ToJson())
}

func (ctx *context) bookmarkDatabasePath() File {
	return NewFileOrPanic(ctx.KlogFolder() + "bookmarks.json")
}

func (ctx *context) Execute(cmd command.Command) Error {
	c := exec.Command(cmd.Bin, cmd.Args...)
	c.Stdin = os.Stdin
	c.Stdout = os.Stdout
	err := c.Run()
	if err == nil {
		return nil
	}
	return NewError(
		"Failed to run command",
		"The command exited with non-zero status",
		err,
	)
}

func (ctx *context) Editors() (string, []command.Command) {
	return ctx.prefs.Editor, POTENTIAL_EDITORS
}

func (ctx *context) FileExplorers() []command.Command {
	return POTENTIAL_FILE_EXLORERS
}

func (ctx *context) Serialiser() parser.Serialiser {
	return ctx.serialiser
}

func (ctx *context) SetSerialiser(s parser.Serialiser) {
	ctx.serialiser = s
}

func (ctx *context) Debug(task func()) {
	if ctx.prefs.IsDebug {
		task()
	}
}

func (ctx *context) Preferences() Preferences {
	return ctx.prefs
}
