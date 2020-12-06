package store

import (
	"github.com/stretchr/testify/assert"
	"klog/datetime"
	"klog/workday"
	"os"
	"testing"
)

const (
	TEST_PATH = "../tmp"
)

func setup() {
	os.RemoveAll(TEST_PATH)
	os.MkdirAll(TEST_PATH, os.ModePerm)
}

func TestInitialisesFileStoreIfPathExists(t *testing.T) {
	setup()
	store, err := CreateFsStore(TEST_PATH)
	assert.Nil(t, err)
	assert.NotNil(t, store)
}

func TestFailsToInitialiseFileStoreIfPathDoesNotExists(t *testing.T) {
	setup()
	store, err := CreateFsStore(TEST_PATH + "/qwerty123")
	assert.Nil(t, store)
	assert.Error(t, err)
}

func TestGetFailsIfDateDoesNotExist(t *testing.T) {
	setup()
	store, _ := CreateFsStore(TEST_PATH)
	_, err := store.Get(datetime.Date{Year: 2020, Month: 1, Day: 31})
	assert.Error(t, err)
}

func TestSavePersists(t *testing.T) {
	setup()
	store, _ := CreateFsStore(TEST_PATH)
	workDay, _ := workday.Create(datetime.Date{Year: 2000, Month: 3, Day: 15})
	err := store.Save(workDay)
	assert.Nil(t, err)
}
