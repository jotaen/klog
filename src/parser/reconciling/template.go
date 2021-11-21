package reconciling

import (
	"errors"
	"github.com/jotaen/klog/src"
	"github.com/jotaen/klog/src/parser"
	"regexp"
	"strings"
	gotime "time"
)

var markerPattern = regexp.MustCompile(`{{.+}}`)

// RenderTemplate replaces placeholders in a template with actual values.
func RenderTemplate(templateText string, time gotime.Time) ([]InsertableText, error) {
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
	_, blocks, err := parser.Parse(instance)
	if err != nil {
		return nil, errors.New("Cannot parse:\n" + instance)
	}
	var texts []InsertableText
	for _, b := range blocks {
		for _, l := range b {
			indentationLevel := 0
			if len(l.PrecedingWhitespace) > 0 {
				indentationLevel = 1
			}
			texts = append(texts, InsertableText{l.Text, indentationLevel})
		}
	}
	return texts, nil
}
