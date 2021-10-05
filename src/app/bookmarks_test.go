package app

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestCreatesNewBookmark(t *testing.T) {
	b := NewBookmark("foo", "/asdf/foo.klg")
	assert.Equal(t, "foo", b.Name().Value())
	assert.Equal(t, "/asdf/foo.klg", b.Target().Path)
}

func TestNormalizesBookmarkName(t *testing.T) {
	b := NewBookmark("@foo", "/asdf/foo.klg")
	assert.Equal(t, "foo", b.Name().Value())

	assert.Equal(t, "foo", NewName("foo").Value())
	assert.Equal(t, "foo", NewName("@foo").Value())
	assert.Equal(t, "foo", NewName("@@foo").Value())
}

func TestCanAddAndRemoveBookmarks(t *testing.T) {
	bc := NewEmptyBookmarksCollection()

	bc.Add(NewDefaultBookmark("/old.klg"))
	assert.Equal(t, "default", bc.Default().Name().Value())
	assert.Equal(t, "/old.klg", bc.Default().Target().Path)
	assert.Equal(t, 1, bc.Count())

	// Overwrites existing bookmark
	bc.Add(NewDefaultBookmark("/new.klg"))
	assert.Equal(t, "/new.klg", bc.Default().Target().Path)
	assert.Equal(t, 1, bc.Count())

	// Add another bookmark
	foo := NewName("foo")
	bc.Add(NewBookmark(foo.Value(), "/qwer.klg"))
	assert.Equal(t, foo, bc.Get(foo).Name())
	assert.Equal(t, 2, bc.Count())

	// Remove again
	bc.Remove(foo)
	assert.Nil(t, bc.Get(foo))
	assert.Equal(t, 1, bc.Count())

	// Clear all
	bc.Clear()
	assert.Nil(t, bc.Default())

	bc.Clear() // Idempotent operation
	assert.Nil(t, bc.Default())
}

func TestParseBookmarksCollectionFromString(t *testing.T) {
	bc, err := NewBookmarksCollectionFromJson(`[{
	"name": "default",
	"path": "/asdf/foo.klg"
}]`)
	require.Nil(t, err)
	def := bc.Default()
	require.NotNil(t, def)
	assert.Equal(t, "default", def.Name().Value())
	assert.Equal(t, "/asdf/foo.klg", def.Target().Path)
}

func TestParseEmptyBookmarksCollectionFromString(t *testing.T) {
	for _, jsonText := range []string{
		``,
		`[]`,
	} {
		bc, err := NewBookmarksCollectionFromJson(jsonText)
		require.Nil(t, err)
		require.NotNil(t, bc)
		assert.Nil(t, bc.Default())
	}
}

func TestParsingFailsForMalformedJson(t *testing.T) {
	for _, json := range []string{
		`[{"name": "defau`, // Invalid JSON
		`{"name": "default", "path": "/asdf/foo.klg"}`, // No array
		`[{"name": "default"}]`,                        // Missing field
		`[{"name": "default", "path": true}]`,          // Wrong type
	} {
		bc, err := NewBookmarksCollectionFromJson(json)
		require.Nil(t, bc)
		assert.Error(t, err)
		assert.Equal(t, CONFIG_ERROR, err.Code())
	}
}

func TestSerializeCollectionToJson(t *testing.T) {
	jsonText := `[
  {
    "name": "default",
    "path": "/asdf.klg"
  }
]
`
	bc, _ := NewBookmarksCollectionFromJson(jsonText)
	assert.Equal(t, jsonText, bc.ToJson())
}

func TestSerializeEmptyCollectionToJson(t *testing.T) {
	jsonText := ``
	bc, _ := NewBookmarksCollectionFromJson(jsonText)
	assert.Equal(t, jsonText, bc.ToJson())
}
