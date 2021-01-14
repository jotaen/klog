package project

import (
	"errors"
	"io/ioutil"
	"klog/datetime"
	parser2 "klog/parser"
	"klog/record"
	"regexp"
	"sort"
	"strconv"
)

type Project interface {
	Name() string
	Path() string
	Get(datetime.Date) (record.Record, []error)
	Save(record.Record) error
	List() ([]datetime.Date, error)
	GetFileProps(record.Record) FileProps
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

func (p project) Get(date datetime.Date) (record.Record, []error) {
	props := createFileProps(p.basePath, date)
	contents, err := readFile(props)
	if err != nil {
		return nil, []error{err}
	}
	workDay, parserErrors := parser2.Parse(contents)
	if parserErrors != nil {
		return nil, parser2.ToErrors(parserErrors)
	}
	return workDay, nil
}

func (p project) Save(workDay record.Record) error {
	props := createFileProps(p.basePath, workDay.Date())
	writeFile(props, parser2.Serialise(workDay))
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
	sort.SliceStable(result, func(i, j int) bool {
		return result[i].ToString() > result[j].ToString()
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

func (p project) GetFileProps(workDay record.Record) FileProps {
	return createFileProps(p.basePath, workDay.Date())
}
