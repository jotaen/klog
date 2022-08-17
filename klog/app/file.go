package app

import (
	"io"
	"os"
	"path/filepath"
)

// File is a descriptor for a file.
// The file is not guaranteed to actually exist on disk.
type File interface {
	// Name returns the file name.
	Name() string

	// Location returns the path of the folder, where the file resides.
	Location() string

	// Path returns the path of the file.
	Path() string
}

// FileWithContents is like File, but with the file contents attached to it.
type FileWithContents interface {
	File

	// Contents returns the contents of the file.
	Contents() string
}

// NewFile creates a new File object from an absolute or relative path.
// It returns an error if the given path cannot be resolved.
func NewFile(path string) (File, Error) {
	absolutePath, err := filepath.Abs(path)
	if err != nil {
		return nil, NewErrorWithCode(
			IO_ERROR,
			"Invalid file path",
			"Location: "+path,
			err,
		)
	}
	return NewFileOrPanic(absolutePath), nil
}

// NewFileOrPanic creates a new File object. It panics, if the path is not absolute.
func NewFileOrPanic(absolutePath string) File {
	if !filepath.IsAbs(absolutePath) {
		panic("Not an absolute path: " + absolutePath)
	}
	return &fileWithPath{absolutePath}
}

type fileWithPath struct {
	absolutePath string
}

type fileWithContents struct {
	File
	contents string
}

func (f *fileWithPath) Name() string {
	return filepath.Base(f.absolutePath)
}

func (f *fileWithPath) Location() string {
	return filepath.Dir(f.absolutePath)
}

func (f *fileWithPath) Path() string {
	return f.absolutePath
}

func (f *fileWithContents) Contents() string {
	return f.contents
}

// IsAbs checks whether the given path is absolute.
func IsAbs(path string) bool {
	return filepath.IsAbs(path)
}

// ReadFile reads the contents of a file from disk and returns it as string.
// It returns an error if the file doesn’t exit or cannot be read.
func ReadFile(source File) (string, Error) {
	contents, err := os.ReadFile(source.Path())
	if err != nil {
		if os.IsNotExist(err) {
			return "", NewErrorWithCode(
				NO_SUCH_FILE,
				"No such file",
				"Location: "+source.Path(),
				err,
			)
		}
		return "", NewErrorWithCode(
			IO_ERROR,
			"Cannot read file",
			"Location: "+source.Path(),
			err,
		)
	}
	return string(contents), nil
}

// WriteToFile saves contents in a file on disk.
// It returns an error if the file cannot be written.
func WriteToFile(target File, contents string) Error {
	err := os.WriteFile(target.Path(), []byte(contents), 0644)
	if err != nil {
		return NewErrorWithCode(
			IO_ERROR,
			"Cannot write to file",
			"Location: "+target.Path(),
			err,
		)
	}
	return nil
}

// ReadStdin reads the entire input from stdin and returns it as string.
// It returns an error if stdin cannot be accessed, or if reading from it fails.
func ReadStdin() (string, Error) {
	stat, err := os.Stdin.Stat()
	if err != nil {
		return "", NewErrorWithCode(
			IO_ERROR,
			"Cannot read from Stdin",
			"Cannot open Stdin stream to check for input",
			err,
		)
	}
	if (stat.Mode() & os.ModeCharDevice) != 0 {
		return "", nil
	}
	bytes, err := io.ReadAll(os.Stdin)
	if err != nil {
		return "", NewErrorWithCode(
			IO_ERROR,
			"Error while reading from Stdin",
			"An error occurred while processing the input stream",
			err,
		)
	}
	return string(bytes), nil
}
