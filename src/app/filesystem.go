package app

import (
	"io/ioutil"
	"os"
)

func fileExists(path string) bool {
	file, err := os.Stat(path)
	if err == nil && file.Mode().IsRegular() {
		return true
	}
	return false
}

func readFile(path string) (string, error) {
	contents, err := ioutil.ReadFile(path)
	if err != nil {
		return "", err
	}
	return string(contents), nil
}

func writeFile(path string, contents string) error {
	//os.MkdirAll(props.Dir, os.ModePerm)
	ioutil.WriteFile(path, []byte(contents), 0644)
	return nil
}
