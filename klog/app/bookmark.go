package app

import (
	"bytes"
	"encoding/json"
	"sort"
	"strings"
)

// Name is the bookmark alias.
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

// Value returns the name of the bookmark without prefix.
func (n Name) Value() string {
	return string(n)
}

// ValuePretty returns the name of the bookmark with prefix.
func (n Name) ValuePretty() string {
	return BOOKMARK_PREFIX + n.Value()
}

// IsValidBookmarkName checks whether `value` is a valid bookmark name (including the prefix).
func IsValidBookmarkName(value string) bool {
	return strings.HasPrefix(value, BOOKMARK_PREFIX)
}

// Bookmark is a way to reference often used files via a short alias (the name).
type Bookmark interface {
	// Name is the alias of the bookmark.
	Name() Name

	// Target is the file that the bookmark references.
	Target() File

	// IsDefault returns whether the bookmark is the default one.
	// In this case, the bookmark name is `default`.
	IsDefault() bool
}

// BookmarksCollection is the collection of all bookmarks.
type BookmarksCollection interface {
	// Get looks up a bookmark by its name.
	Get(Name) Bookmark

	// All returns all bookmarks in the collection.
	All() []Bookmark

	// Default returns the default bookmark of the collection.
	Default() Bookmark

	// Set adds a new bookmark to the collection.
	Set(Bookmark)

	// Remove deletes a bookmark from the collection.
	Remove(Name) bool

	// Clear deletes all bookmarks of the collection.
	Clear()

	// ToJson returns a JSON-representation of the bookmark collection.
	ToJson() string

	// Count returns the number of bookmarks in the collection.
	Count() int
}

func NewBookmark(name string, target File) Bookmark {
	return &bookmark{NewName(name), target}
}

func NewDefaultBookmark(target File) Bookmark {
	return NewBookmark(BOOKMARK_DEFAULT_NAME, target)
}

type bookmark struct {
	name   Name
	target File
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

// NewBookmarksCollectionFromJson deserialises JSON data. It returns an error
// if the syntax is malformed.
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
		if !IsAbs(*b.Path) {
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
