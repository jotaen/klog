package store

import (
	"errors"
	"github.com/stretchr/testify/assert"
	"klog/datetime"
	"klog/testutil"
	"klog/workday"
	"testing"
)

func TestInitialisesFileStoreIfPathExists(t *testing.T) {
	testutil.WithDisk(func(path string) {
		store, err := CreateFsStore(path)
		assert.Nil(t, err)
		assert.NotNil(t, store)
	})
}

func TestFailsToInitialiseFileStoreIfPathDoesNotExists(t *testing.T) {
	testutil.WithDisk(func(path string) {
		store, err := CreateFsStore(path + "/qwerty123")
		assert.Nil(t, store)
		assert.Equal(t, err, errors.New("NO_SUCH_PATH"))
	})
}

func TestGetFailsIfDateDoesNotExist(t *testing.T) {
	testutil.WithDisk(func(path string) {
		store, _ := CreateFsStore(path)
		date, _ := datetime.CreateDate(2020, 1, 31)
		_, errs := store.Get(date)
		assert.Error(t, errs[0])
	})
}

func TestSavePersists(t *testing.T) {
	testutil.WithDisk(func(path string) {
		store, _ := CreateFsStore(path)
		date, _ := datetime.CreateDate(1999, 3, 15)
		originalWd := workday.Create(date)
		err := store.Save(originalWd)
		assert.Nil(t, err)

		readWd, _ := store.Get(date)
		assert.Equal(t, originalWd, readWd)
	})
}

func TestListReturnsPersistedWorkdays(t *testing.T) {
	testutil.WithDisk(func(path string) {
		store, _ := CreateFsStore(path)

		date1, _ := datetime.CreateDate(1999, 1, 13)
		store.Save(workday.Create(date1))
		date2, _ := datetime.CreateDate(1999, 1, 14)
		store.Save(workday.Create(date2))
		date3, _ := datetime.CreateDate(1999, 2, 5)
		store.Save(workday.Create(date3))

		wds, _ := store.List()
		assert.Equal(t, []datetime.Date{date1, date2, date3}, wds)
	})
}
