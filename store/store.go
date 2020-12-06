package store

import (
	"errors"
	"fmt"
	"klog/workday"
	"os"
)

type Store interface {
	Get(workday.Date) (workday.WorkDay, error)
	// List() ([]workday.WorkDay, error)
	// Save(workday.WorkDay) error
}

type fileStore struct {
	basePath string
}

func CreateFsStore(path string) (Store, error) {
	if !dirExists(path) {
		return nil, errors.New("Not such directory")
	}
	return fileStore{
		basePath: path,
	}, nil
}

func (fs fileStore) Get(date workday.Date) (workday.WorkDay, error) {
	if fileExists(fs.makePath(date)) {
		return nil, nil
	}
	return nil, errors.New("No such entry")
}

func (fs fileStore) makePath(date workday.Date) string {
	return fmt.Sprintf("%v/%v/%v/%v", fs.basePath, date.Year, date.Month, date.Day)
}

func dirExists(path string) bool {
	file, err := os.Stat(path)
	if err == nil && file.Mode().IsDir() {
		return true
	}
	return false
}

func fileExists(path string) bool {
	file, err := os.Stat(path)
	if err == nil && file.Mode().IsRegular() {
		return true
	}
	return false
}
