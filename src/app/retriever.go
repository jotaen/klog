package app

import (
	"errors"
	"strings"
)

type Retriever func(fileArgs ...FileOrBookmarkName) ([]*fileWithContent, Error)

type fileRetriever struct {
	readFile  func(File) (string, Error)
	bookmarks BookmarksCollection
}

type fileWithContent struct {
	File
	content string
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

func (retriever *fileRetriever) Retrieve(fileArgs ...FileOrBookmarkName) ([]*fileWithContent, Error) {
	fileArgs = removeBlankEntries(fileArgs...)
	if len(fileArgs) == 0 {
		defaultBookmark := retriever.bookmarks.Default()
		if defaultBookmark != nil {
			fileArgs = []FileOrBookmarkName{
				FileOrBookmarkName(defaultBookmark.Target().Path()),
			}
		}
	}
	var results []*fileWithContent
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
		results = append(results, &fileWithContent{file, content})
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

type stdinRetriever struct {
	readStdin func() (string, Error)
}

func (retriever *stdinRetriever) Retrieve(fileArgs ...FileOrBookmarkName) ([]*fileWithContent, Error) {
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
	return []*fileWithContent{{
		File:    nil,
		content: stdin,
	}}, nil
}

func retrieveFirst(rs []Retriever, fileArgs ...FileOrBookmarkName) ([]*fileWithContent, Error) {
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
