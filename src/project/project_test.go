package project

import (
	"errors"
	"github.com/stretchr/testify/assert"
	"klog/datetime"
	datetime2 "klog/testutil/datetime"
	. "klog/testutil/withdisk"
	"klog/record"
	"testing"
)

func TestInitialisesFileStoreIfPathExists(t *testing.T) {
	WithDisk(func(path string) {
		store, err := NewProject(path)
		assert.Nil(t, err)
		assert.NotNil(t, store)
	})
}

func TestFailsToInitialiseFileStoreIfPathDoesNotExists(t *testing.T) {
	WithDisk(func(path string) {
		store, err := NewProject(path + "/qwerty123")
		assert.Nil(t, store)
		assert.Equal(t, err, errors.New("NO_SUCH_PATH"))
	})
}

func TestGetFailsIfDateDoesNotExist(t *testing.T) {
	WithDisk(func(path string) {
		store, _ := NewProject(path)
		_, errs := store.Get(datetime2.Date_(2020, 1, 31))
		assert.Error(t, errs[0])
	})
}

func TestSavePersists(t *testing.T) {
	WithDisk(func(path string) {
		store, _ := NewProject(path)
		date := datetime2.Date_(1999, 3, 15)
		originalWd := record.NewRecord(date)
		err := store.Save(originalWd)
		assert.Nil(t, err)

		readWd, _ := store.Get(date)
		assert.Equal(t, originalWd, readWd)
	})
}

func TestListReturnsPersistedWorkdaysInDescendingOrder(t *testing.T) {
	WithDisk(func(path string) {
		store, _ := NewProject(path)

		date1 := datetime2.Date_(1999, 2, 5)
		store.Save(record.NewRecord(date1))
		date2 := datetime2.Date_(1999, 1, 14)
		store.Save(record.NewRecord(date2))
		date3 := datetime2.Date_(1999, 1, 13)
		store.Save(record.NewRecord(date3))

		wds, _ := store.List()
		assert.Equal(t, []datetime.Date{date1, date2, date3}, wds)
	})
}

func TestListReturnsFilteredWorkdays(t *testing.T) {
	WithDisk(func(path string) {
		store, _ := NewProject(path)

		date1 := datetime2.Date_(1999, 2, 5)
		store.Save(record.NewRecord(date1))
		date2 := datetime2.Date_(1999, 1, 14)
		store.Save(record.NewRecord(date2))
		date3 := datetime2.Date_(1999, 1, 13)
		store.Save(record.NewRecord(date3))

		wds, _ := store.List()
		assert.Equal(t, []datetime.Date{date1, date2, date3}, wds)
	})
}
