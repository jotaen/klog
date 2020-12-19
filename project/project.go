package project

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

type Project interface {
	Name() string
	Path() string
	Get(datetime.Date) (workday.WorkDay, []error)
	Save(workday.WorkDay) error
	List() ([]datetime.Date, error)
	GetFileProps(workday.WorkDay) FileProps
}

type project struct {
	basePath string
}

func NewProject(path string) (Project, error) {
	if !dirExists(path) {
		return nil, errors.New("NO_SUCH_PATH")
	}
	return project{
		basePath: path,
	}, nil
}

func (p project) Name() string {
	return p.basePath // TODO
}

func (p project) Path() string {
	return p.basePath
}

func (p project) Get(date datetime.Date) (workday.WorkDay, []error) {
	props := createFileProps(p.basePath, date)
	contents, err := readFile(props)
	if err != nil {
		return nil, []error{err}
	}
	workDay, parserErrors := parser.Parse(contents)
	if parserErrors != nil {
		return nil, parser.ToErrors(parserErrors)
	}
	return workDay, nil
}

func (p project) Save(workDay workday.WorkDay) error {
	props := createFileProps(p.basePath, workDay.Date())
	writeFile(props, serialiser.Serialise(workDay))
	return nil
}

var datePattern = regexp.MustCompile("^[0-9]{4}$")
var monthPattern = regexp.MustCompile("^[0-9]{2}$")
var dayPattern = regexp.MustCompile("^[0-9]{2}.yml$")

func (p project) List() ([]datetime.Date, error) {
	var result []datetime.Date
	walkDir(p.basePath, true, datePattern, func(year string) {
		walkDir(p.basePath+"/"+year, true, monthPattern, func(month string) {
			walkDir(p.basePath+"/"+year+"/"+month, false, dayPattern, func(day string) {
				yyyy, _ := strconv.Atoi(year)
				mm, _ := strconv.Atoi(month)
				dd, _ := strconv.Atoi(day[0:2])
				date, err := datetime.NewDate(yyyy, mm, dd)
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

func (p project) GetFileProps(workDay workday.WorkDay) FileProps {
	return createFileProps(p.basePath, workDay.Date())
}
