package app

import (
	"errors"
	"strings"
)

type Retriever func(fileArgs ...string) ([]fileWithContent, Error)

type fileRetriever struct {
	readFile  func(string) (string, Error)
	bookmarks BookmarksCollection
}

type fileWithContent struct {
	File
	content string
}

func removeBlankEntries(fileArgs ...string) []string {
	var result []string
	for _, f := range fileArgs {
		if strings.TrimLeft(f, " ") == "" {
			continue
		}
		result = append(result, f)
	}
	return result
}

func (ir *fileRetriever) Retrieve(fileArgs ...string) ([]fileWithContent, Error) {
	fileArgs = removeBlankEntries(fileArgs...)
	if len(fileArgs) == 0 {
		defaultBookmark := ir.bookmarks.Default()
		if defaultBookmark != nil {
			fileArgs = []string{defaultBookmark.Target().Path()}
		}
	}
	var results []fileWithContent
	var errs []string
	for _, arg := range fileArgs {
		path, pathErr := (func() (string, error) {
			if strings.HasPrefix(arg, "@") {
				b := ir.bookmarks.Get(NewName(arg))
				if b == nil {
					return arg, errors.New("No such bookmark")
				}
				return b.Target().Path(), nil
			}
			return arg, nil
		})()
		if pathErr != nil {
			errs = append(errs, pathErr.Error()+": "+arg)
			continue
		}
		content, readErr := ir.readFile(path)
		if readErr != nil {
			errs = append(errs, readErr.Error()+": "+path)
			continue
		}
		results = append(results, fileWithContent{NewFile(path), content})
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

func (r *stdinRetriever) Retrieve(fileArgs ...string) ([]fileWithContent, Error) {
	fileArgs = removeBlankEntries(fileArgs...)
	if len(fileArgs) > 0 {
		return nil, nil
	}
	stdin, err := r.readStdin()
	if err != nil {
		return nil, err
	}
	if stdin == "" {
		return nil, nil
	}
	return []fileWithContent{{
		File:    NewFile("/dev/stdin"), // Fake file just to fulfill interface
		content: stdin,
	}}, nil
}

func retrieveFirst(rs []Retriever, fileArgs ...string) ([]fileWithContent, Error) {
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
