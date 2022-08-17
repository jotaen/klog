package app

import (
	"errors"
	"strings"
)

// Retriever is the function interface for retrieving the input data from
// various kinds of sources.
type Retriever func(fileArgs ...FileOrBookmarkName) ([]FileWithContents, Error)

type FileRetriever struct {
	readFile  func(File) (string, Error)
	bookmarks BookmarksCollection
}

// Retrieve retrieves the contents from files or bookmarks. If no arguments were
// specified, it tries to read from the default bookmark.
func (retriever *FileRetriever) Retrieve(fileArgs ...FileOrBookmarkName) ([]FileWithContents, Error) {
	fileArgs = removeBlankEntries(fileArgs...)
	if len(fileArgs) == 0 {
		defaultBookmark := retriever.bookmarks.Default()
		if defaultBookmark != nil {
			fileArgs = []FileOrBookmarkName{
				FileOrBookmarkName(defaultBookmark.Target().Path()),
			}
		}
	}
	var results []FileWithContents
	var errs []string
	for _, arg := range fileArgs {
		argValue := string(arg)
		path, pathErr := (func() (string, error) {
			if IsValidBookmarkName(argValue) {
				b := retriever.bookmarks.Get(NewName(argValue))
				if b == nil {
					return argValue, errors.New("No such bookmark")
				}
				return b.Target().Path(), nil
			}
			return argValue, nil
		})()
		if pathErr != nil {
			errs = append(errs, pathErr.Error()+": "+argValue)
			continue
		}
		file, fErr := NewFile(path)
		if fErr != nil {
			errs = append(errs, "Invalid file path: "+path)
		}
		content, readErr := retriever.readFile(file)
		if readErr != nil {
			errs = append(errs, readErr.Error()+": "+file.Path())
			continue
		}
		results = append(results, &fileWithContents{file, content})
	}
	if len(errs) > 0 {
		return nil, NewErrorWithCode(
			IO_ERROR,
			"Cannot retrieve files",
			strings.Join(errs, "\n"),
			nil,
		)
	}
	return results, nil
}

type StdinRetriever struct {
	readStdin func() (string, Error)
}

// Retrieve retrieves the content from stdin. It only returns something if no
// file or bookmark names were specified explicitly.
func (retriever *StdinRetriever) Retrieve(fileArgs ...FileOrBookmarkName) ([]FileWithContents, Error) {
	fileArgs = removeBlankEntries(fileArgs...)
	if len(fileArgs) > 0 {
		return nil, nil
	}
	stdin, err := retriever.readStdin()
	if err != nil {
		return nil, err
	}
	if stdin == "" {
		return nil, nil
	}
	return []FileWithContents{&fileWithContents{
		File:     &fileWithPath{""},
		contents: stdin,
	}}, nil
}

func removeBlankEntries(fileArgs ...FileOrBookmarkName) []FileOrBookmarkName {
	var result []FileOrBookmarkName
	for _, f := range fileArgs {
		if strings.TrimLeft(string(f), " ") == "" {
			continue
		}
		result = append(result, f)
	}
	return result
}

// retrieveFirst returns the result from the first Retriever that yields something.
func retrieveFirst(rs []Retriever, fileArgs ...FileOrBookmarkName) ([]FileWithContents, Error) {
	for _, r := range rs {
		files, err := r(fileArgs...)
		if err != nil {
			return nil, err
		}
		if len(files) > 0 {
			return files, nil
		}
	}
	return nil, nil
}
