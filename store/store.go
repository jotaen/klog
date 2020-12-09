package store

import (
	"errors"
	"io/ioutil"
	"klog/datetime"
	"klog/parser"
	"klog/serialiser"
	"klog/workday"
	"regexp"
	"strconv"
)

type Store interface {
	Get(datetime.Date) (workday.WorkDay, []error)
	Save(workday.WorkDay) error
	List() ([]datetime.Date, error)
	GetFileProps(workday.WorkDay) FileProps
}

type fileStore struct {
	basePath string
}

func CreateFsStore(path string) (Store, error) {
	if !dirExists(path) {
		return nil, errors.New("NO_SUCH_PATH")
	}
	return fileStore{
		basePath: path,
	}, nil
}

func (fs fileStore) Get(date datetime.Date) (workday.WorkDay, []error) {
	props := createFileProps(fs.basePath, date)
	contents, err := readFile(props)
	if err != nil {
		return nil, []error{err}
	}
	workDay, errs := parser.Parse(contents)
	return workDay, errs
}

func (fs fileStore) Save(workDay workday.WorkDay) error {
	props := createFileProps(fs.basePath, workDay.Date())
	writeFile(props, serialiser.Serialise(workDay))
	return nil
}

var datePattern = regexp.MustCompile("^[0-9]{4}$")
var monthPattern = regexp.MustCompile("^[0-9]{2}$")
var dayPattern = regexp.MustCompile("^[0-9]{2}.yml$")

func (fs fileStore) List() ([]datetime.Date, error) {
	result := []datetime.Date{}
	walkDir(fs.basePath, true, datePattern, func(year string) {
		walkDir(fs.basePath+"/"+year, true, monthPattern, func(month string) {
			walkDir(fs.basePath+"/"+year+"/"+month, false, dayPattern, func(day string) {
				yyyy, _ := strconv.Atoi(year)
				mm, _ := strconv.Atoi(month)
				dd, _ := strconv.Atoi(day[0:2])
				date, err := datetime.CreateDate(yyyy, mm, dd)
				if err == nil {
					result = append(result, date)
				}
			})
		})
	})
	return result, nil
}

func walkDir(
	path string,
	mustBeDir bool,
	pattern *regexp.Regexp,
	fn func(string),
) []string {
	files, err := ioutil.ReadDir(path)
	result := []string{}
	if err != nil {
		return result
	}
	for _, file := range files {
		if (mustBeDir == file.IsDir()) && pattern.MatchString(file.Name()) {
			fn(file.Name())
		}
	}
	return result
}

func (fs fileStore) GetFileProps(workDay workday.WorkDay) FileProps {
	return createFileProps(fs.basePath, workDay.Date())
}
