package store

import (
	"errors"
	"github.com/stretchr/testify/assert"
	"klog/datetime"
	datetime2 "klog/testutil/datetime"
	. "klog/testutil/withdisk"
	"klog/workday"
	"testing"
)

func TestInitialisesFileStoreIfPathExists(t *testing.T) {
	WithDisk(func(path string) {
		store, err := NewFsStore(path)
		assert.Nil(t, err)
		assert.NotNil(t, store)
	})
}

func TestFailsToInitialiseFileStoreIfPathDoesNotExists(t *testing.T) {
	WithDisk(func(path string) {
		store, err := NewFsStore(path + "/qwerty123")
		assert.Nil(t, store)
		assert.Equal(t, err, errors.New("NO_SUCH_PATH"))
	})
}

func TestGetFailsIfDateDoesNotExist(t *testing.T) {
	WithDisk(func(path string) {
		store, _ := NewFsStore(path)
		_, errs := store.Get(datetime2.Date_(2020, 1, 31))
		assert.Error(t, errs[0])
	})
}

func TestSavePersists(t *testing.T) {
	WithDisk(func(path string) {
		store, _ := NewFsStore(path)
		date := datetime2.Date_(1999, 3, 15)
		originalWd := workday.NewWorkDay(date)
		err := store.Save(originalWd)
		assert.Nil(t, err)

		readWd, _ := store.Get(date)
		assert.Equal(t, originalWd, readWd)
	})
}

func TestListReturnsPersistedWorkdays(t *testing.T) {
	WithDisk(func(path string) {
		store, _ := NewFsStore(path)

		date1 := datetime2.Date_(1999, 1, 13)
		store.Save(workday.NewWorkDay(date1))
		date2 := datetime2.Date_(1999, 1, 14)
		store.Save(workday.NewWorkDay(date2))
		date3 := datetime2.Date_(1999, 2, 5)
		store.Save(workday.NewWorkDay(date3))

		wds, _ := store.List()
		assert.Equal(t, []datetime.Date{date1, date2, date3}, wds)
	})
}
