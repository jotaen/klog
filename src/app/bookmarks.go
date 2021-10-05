package app

import (
	"bytes"
	"encoding/json"
	"strings"
)

type Name string

var defaultName = "default"

func NewName(name string) Name {
	return Name(strings.TrimLeft(name, "@"))
}

func (n Name) Value() string {
	return string(n)
}

type Bookmark interface {
	Name() Name
	Target() *File
	IsDefault() bool
}

type bookmark struct {
	name   Name
	target *File
}

func NewBookmark(name string, targetPath string) Bookmark {
	return &bookmark{NewName(name), NewFile(targetPath)}
}

func NewDefaultBookmark(targetPath string) Bookmark {
	return NewBookmark(defaultName, targetPath)
}

func (b *bookmark) Name() Name {
	return b.name
}

func (b *bookmark) Target() *File {
	return b.target
}

func (b *bookmark) IsDefault() bool {
	return b.name.Value() == defaultName
}

type BookmarksCollection interface {
	Get(Name) Bookmark
	Default() Bookmark
	Add(Bookmark)
	Remove(Name)
	Clear()
	Count() int
	ToJson() string
}

type bookmarksCollection struct {
	bookmarks map[Name]Bookmark
}

func (bc *bookmarksCollection) Default() Bookmark {
	return bc.bookmarks[Name(defaultName)]
}

type bookmarkJson struct {
	Name *string `json:"name"`
	Path *string `json:"path"`
}

func NewEmptyBookmarksCollection() BookmarksCollection {
	return &bookmarksCollection{make(map[Name]Bookmark)}
}

func NewBookmarksCollectionFromJson(jsonText string) (BookmarksCollection, Error) {
	newMalformedJsonError := func(err error) Error {
		return NewErrorWithCode(
			CONFIG_ERROR,
			"Invalid JSON",
			"The JSON in your bookmarks file is malformed",
			err,
		)
	}
	bc := NewEmptyBookmarksCollection()
	if jsonText == "" {
		return bc, nil
	}
	var rawBookmarkInfo []bookmarkJson
	err := json.Unmarshal([]byte(jsonText), &rawBookmarkInfo)
	if err != nil {
		return nil, newMalformedJsonError(err)
	}
	for _, b := range rawBookmarkInfo {
		if b.Name == nil || b.Path == nil {
			return nil, newMalformedJsonError(nil)
		}
		bc.Add(NewBookmark(*b.Name, *b.Path))
	}
	return bc, nil
}

func (bc *bookmarksCollection) Get(n Name) Bookmark {
	return bc.bookmarks[n]
}

func (bc *bookmarksCollection) Add(b Bookmark) {
	bc.bookmarks[b.Name()] = b
}

func (bc *bookmarksCollection) Remove(n Name) {
	delete(bc.bookmarks, n)
}

func (bc *bookmarksCollection) Clear() {
	bc.bookmarks = make(map[Name]Bookmark)
}

func (bc *bookmarksCollection) Count() int {
	return len(bc.bookmarks)
}

func (bc *bookmarksCollection) ToJson() string {
	if bc.Default() == nil {
		return ""
	}
	name := bc.Default().Name().Value()
	path := bc.Default().Target().Path
	bookmarksAsJson := []bookmarkJson{
		{&name, &path},
	}
	buffer := new(bytes.Buffer)
	enc := json.NewEncoder(buffer)
	enc.SetIndent("", "  ")
	enc.SetEscapeHTML(false)
	err := enc.Encode(&bookmarksAsJson)
	if err != nil {
		panic(err)
	}
	return buffer.String()
}
