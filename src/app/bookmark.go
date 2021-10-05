package app

import (
	"bytes"
	"encoding/json"
	"path"
	"sort"
	"strings"
)

type Name string

const (
	BOOKMARK_DEFAULT_NAME = "default"
	BOOKMARK_PREFIX       = "@"
)

func NewName(name string) Name {
	value := strings.TrimLeft(name, BOOKMARK_PREFIX)
	if value == "" {
		value = BOOKMARK_DEFAULT_NAME
	}
	return Name(value)
}

func (n Name) Value() string {
	return string(n)
}

func (n Name) ValuePretty() string {
	return BOOKMARK_PREFIX + n.Value()
}

func IsValidBookmarkName(value string) bool {
	return strings.HasPrefix(value, BOOKMARK_PREFIX)
}

type Bookmark interface {
	Name() Name
	Target() File
	IsDefault() bool
}

type bookmark struct {
	name   Name
	target File
}

func NewBookmark(name string, target File) Bookmark {
	return &bookmark{NewName(name), target}
}

func NewDefaultBookmark(target File) Bookmark {
	return NewBookmark(BOOKMARK_DEFAULT_NAME, target)
}

func (b *bookmark) Name() Name {
	return b.name
}

func (b *bookmark) Target() File {
	return b.target
}

func (b *bookmark) IsDefault() bool {
	return b.name.Value() == BOOKMARK_DEFAULT_NAME
}

type BookmarksCollection interface {
	Get(Name) Bookmark
	All() []Bookmark
	Default() Bookmark
	Set(Bookmark)
	Remove(Name) bool
	Clear()
	ToJson() string
	Count() int
}

type bookmarksCollection struct {
	bookmarks map[Name]Bookmark
}

func (bc *bookmarksCollection) Default() Bookmark {
	return bc.bookmarks[Name(BOOKMARK_DEFAULT_NAME)]
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
		if !path.IsAbs(*b.Path) {
			return nil, newMalformedJsonError(nil)
		}
		file, fErr := NewFile(*b.Path)
		if fErr != nil {
			return nil, fErr
		}
		bc.Set(NewBookmark(*b.Name, file))
	}
	return bc, nil
}

func (bc *bookmarksCollection) Get(n Name) Bookmark {
	return bc.bookmarks[n]
}

func (bc *bookmarksCollection) All() []Bookmark {
	sortedBookmarks := make([]Bookmark, 0, len(bc.bookmarks))
	for _, b := range bc.bookmarks {
		sortedBookmarks = append(sortedBookmarks, b)
	}
	sort.Slice(sortedBookmarks, func(i, j int) bool {
		return sortedBookmarks[i].Name() < sortedBookmarks[j].Name()
	})
	return sortedBookmarks
}

func (bc *bookmarksCollection) Set(b Bookmark) {
	bc.bookmarks[b.Name()] = b
}

func (bc *bookmarksCollection) Remove(n Name) bool {
	if bc.bookmarks[n] == nil {
		return false
	}
	delete(bc.bookmarks, n)
	return true
}

func (bc *bookmarksCollection) Clear() {
	bc.bookmarks = make(map[Name]Bookmark)
}

func (bc *bookmarksCollection) Count() int {
	return len(bc.bookmarks)
}

func (bc *bookmarksCollection) ToJson() string {
	var bookmarksAsJson []bookmarkJson
	for _, b := range bc.All() {
		name := b.Name().Value()
		path := b.Target().Path()
		bookmarksAsJson = append(bookmarksAsJson, bookmarkJson{
			&name, &path,
		})
	}
	if len(bookmarksAsJson) == 0 {
		return ""
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
