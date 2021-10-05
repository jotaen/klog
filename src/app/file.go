package app

import (
	"io"
	"os"
	"path"
	"path/filepath"
)

type File interface {
	Name() string
	Location() string
	Path() string
}

type absoluteFile struct {
	absolute string
}

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

func NewFileOrPanic(absolutePath string) File {
	if !path.IsAbs(absolutePath) {
		panic("Not an absolute path: " + absolutePath)
	}
	return &absoluteFile{absolutePath}
}

func (f *absoluteFile) Name() string {
	return filepath.Base(f.absolute)
}

func (f *absoluteFile) Location() string {
	return filepath.Dir(f.absolute)
}

func (f *absoluteFile) Path() string {
	return f.absolute
}

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
