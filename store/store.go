package store

import (
	"errors"
	"fmt"
	"klog/datetime"
	"klog/parser"
	"klog/serialiser"
	"klog/workday"
)

type Store interface {
	Get(datetime.Date) (workday.WorkDay, []error)
	Save(workday.WorkDay) error
	// List() ([]workday.WorkDay, error)
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

func (fs fileStore) Get(date datetime.Date) (workday.WorkDay, []error) {
	props := fs.createFileProps(date)
	contents, err := readFile(props)
	if err != nil {
		return nil, []error{err}
	}
	workDay, errs := parser.Parse(contents)
	return workDay, errs
}

func (fs fileStore) Save(workDay workday.WorkDay) error {
	props := fs.createFileProps(workDay.Date())
	writeFile(props, serialiser.Serialise(workDay))
	return nil
}

func (fs fileStore) createFileProps(date datetime.Date) fileProps {
	props := fileProps{
		dir:  fmt.Sprintf("%v/%v/%02v", fs.basePath, date.Year, date.Month),
		name: fmt.Sprintf("%02v.yml", date.Day),
	}
	props.path = props.dir + "/" + props.name
	return props
}
