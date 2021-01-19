package app

import (
	"io/ioutil"
	"os"
	"path/filepath"
)

func fileExists(path string) bool {
	file, err := os.Stat(path)
	if err == nil && file.Mode().IsRegular() {
		return true
	}
	return false
}

func ReadFile(path string) (string, error) {
	contents, err := ioutil.ReadFile(path)
	if err != nil {
		return "", err
	}
	return string(contents), nil
}

func WriteFile(path string, contents string) error {
	dir := filepath.Dir(path)
	err := os.MkdirAll(dir, os.ModePerm)
	if err != nil {
		return err
	}
	return ioutil.WriteFile(path, []byte(contents), 0644)
}
