package project

import (
	"fmt"
	"io/ioutil"
	"klog/datetime"
	"os"
)

type FileProps struct {
	Dir  string
	Name string
	Path string
}

func createFileProps(basePath string, date datetime.Date) FileProps {
	props := FileProps{
		Dir:  fmt.Sprintf("%v/%v/%02v", basePath, date.Year(), date.Month()),
		Name: fmt.Sprintf("%02v.yml", date.Day()),
	}
	props.Path = props.Dir + "/" + props.Name
	return props
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

func readFile(props FileProps) (string, error) {
	contents, err := ioutil.ReadFile(props.Path)
	if err != nil {
		return "", err
	}
	return string(contents), nil
}

func writeFile(props FileProps, contents string) error {
	os.MkdirAll(props.Dir, os.ModePerm)
	ioutil.WriteFile(props.Path, []byte(contents), 0644)
	return nil
}
