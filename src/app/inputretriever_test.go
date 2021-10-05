package app

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

type MockFs map[string]bool

func (fs MockFs) readFile(path string) (string, Error) {
	if fs[path] {
		return path, nil
	}
	return "", NewError("", path, nil)
}

func TestFileRetrieverResolvesFilesAndBookmarks(t *testing.T) {
	bc := NewEmptyBookmarksCollection()
	bc.Add(NewBookmark("foo", "/foo.klg"))
	files, err := (&fileRetriever{
		MockFs{"/asdf.klg": true, "/foo.klg": true}.readFile,
		bc,
	}).Retrieve("/asdf.klg", "@foo")

	require.Nil(t, err)
	require.Len(t, files, 2)
	assert.Equal(t, "/asdf.klg", files[0].Path)
	assert.Equal(t, "/foo.klg", files[1].Path)
}

func TestReturnsErrorIfBookmarksOrFilesAreInvalid(t *testing.T) {
	bc := NewEmptyBookmarksCollection()
	bc.Add(NewBookmark("foo", "/foo.klg"))
	files, err := (&fileRetriever{
		MockFs{}.readFile,
		bc,
	}).Retrieve("/asdf.klg", "@foo", "@bar")

	require.Nil(t, files)
	require.Error(t, err)
	assert.Contains(t, err.Details(), "/asdf.klg")
	assert.Contains(t, err.Details(), "/foo.klg")
	assert.Contains(t, err.Details(), "@bar")
}

func TestFallsBackToDefaultBookmark(t *testing.T) {
	bc := NewEmptyBookmarksCollection()
	bc.Add(NewDefaultBookmark("/foo.klg"))
	retriever := &fileRetriever{
		MockFs{"/foo.klg": true}.readFile,
		bc,
	}
	for _, f := range []func() ([]fileWithContent, Error){
		func() ([]fileWithContent, Error) { return retriever.Retrieve() },
		func() ([]fileWithContent, Error) { return retriever.Retrieve("") },
		func() ([]fileWithContent, Error) { return retriever.Retrieve("", " ") },
	} {
		files, err := f()
		require.Nil(t, err)
		require.Len(t, files, 1)
		assert.Equal(t, "/foo.klg", files[0].Path)
	}
}

func TestReturnsStdinInput(t *testing.T) {
	retriever := &stdinRetriever{
		func() (string, Error) { return "2021-01-01", nil },
	}
	for _, f := range []func() ([]fileWithContent, Error){
		func() ([]fileWithContent, Error) { return retriever.Retrieve() },
		func() ([]fileWithContent, Error) { return retriever.Retrieve("") },
		func() ([]fileWithContent, Error) { return retriever.Retrieve("", " ") },
	} {
		files, err := f()
		require.Nil(t, err)
		require.Len(t, files, 1)
		assert.Equal(t, "2021-01-01", files[0].content)
	}
}

func TestBailsOutIfFileArgsGiven(t *testing.T) {
	files, err := (&stdinRetriever{
		func() (string, Error) { return "", nil },
	}).Retrieve("foo.klg")

	require.Nil(t, err)
	require.Nil(t, files)
}
