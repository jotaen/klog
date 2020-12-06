package store

import (
	"github.com/stretchr/testify/assert"
	"klog/workday"
	"os"
	"testing"
	"time"
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
	_, err := store.Get(workday.Date{Year: 2020, Month: time.January, Day: 31})
	assert.Error(t, err)
}
