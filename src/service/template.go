package service

import (
	"errors"
	"klog"
	"klog/parser"
	"regexp"
	"strings"
	gotime "time"
)

type RecordText string

func RenderTemplate(templateText string, time gotime.Time) (RecordText, error) {
	today := klog.NewDateFromTime(time)
	now := klog.NewTimeFromTime(time)
	variables := map[string]string{
		"TODAY":     today.ToString(),
		"YESTERDAY": today.PlusDays(-1).ToString(),
		"NOW":       now.ToString(),
	}
	instance := markerPattern.ReplaceAllStringFunc(templateText, func(m string) string {
		m = strings.TrimLeft(m, "{{")
		m = strings.TrimRight(m, "}}")
		m = strings.TrimSpace(m)
		return variables[m]
	})
	_, err := parser.Parse(instance)
	if err != nil {
		return "", errors.New("Cannot parse:\n" + instance)
	}
	return RecordText(instance), nil
}

func AppendableText(content string, newRecordText RecordText) string {
	result := string(newRecordText)
	if content == "" {
		return result
	}
	if strings.HasSuffix(content, "\n\n") {
		return result
	}
	if strings.HasSuffix(content, "\n") {
		return "\n" + result
	}
	return "\n\n" + result
}

var markerPattern = regexp.MustCompile(`{{.+}}`)
