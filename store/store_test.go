package store

import (
	"github.com/stretchr/testify/assert"
	"klog/datetime"
	"klog/workday"
	"os"
	"testing"
)

func run(fn func(string)) {
	path := "../tmp"
	os.RemoveAll(path)
	os.MkdirAll(path, os.ModePerm)
	fn(path)
	os.RemoveAll(path)
}

func TestInitialisesFileStoreIfPathExists(t *testing.T) {
	run(func(path string) {
		store, err := CreateFsStore(path)
		assert.Nil(t, err)
		assert.NotNil(t, store)
	})
}

func TestFailsToInitialiseFileStoreIfPathDoesNotExists(t *testing.T) {
	run(func(path string) {
		store, err := CreateFsStore(path + "/qwerty123")
		assert.Nil(t, store)
		assert.Error(t, err)
	})
}

func TestGetFailsIfDateDoesNotExist(t *testing.T) {
	run(func(path string) {
		store, _ := CreateFsStore(path)
		date, _ := datetime.CreateDate(2020, 1, 31)
		_, errs := store.Get(date)
		assert.Error(t, errs[0])
	})
}

func TestSavePersists(t *testing.T) {
	run(func(path string) {
		store, _ := CreateFsStore(path)
		date, _ := datetime.CreateDate(1999, 3, 15)
		originalWd := workday.Create(date)
		err := store.Save(originalWd)
		assert.Nil(t, err)

		readWd, _ := store.Get(date)
		assert.Equal(t, originalWd, readWd)
	})
}
