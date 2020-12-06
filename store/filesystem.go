package store

import (
	"io/ioutil"
	"os"
)

type fileProps struct {
	dir  string
	name string
	path string
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

func readFile(props fileProps) (string, error) {
	contents, err := ioutil.ReadFile(props.path)
	if err != nil {
		return "", err
	}
	return string(contents), nil
}

func writeFile(props fileProps, contents string) error {
	os.MkdirAll(props.dir, os.ModePerm)
	ioutil.WriteFile(props.path, []byte(contents), 0644)
	return nil
}
