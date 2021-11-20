package reconciler

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
func RenderTemplate(templateText string, time gotime.Time) ([]Text, error) {
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
	pr, err := parser.Parse(instance)
	if err != nil {
		return nil, errors.New("Cannot parse:\n" + instance)
	}
	var texts []Text
	for _, l := range pr.Lines {
		indentationLevel := 0
		if len(l.PrecedingWhitespace) > 0 {
			indentationLevel = 1
		}
		texts = append(texts, Text{l.Text, indentationLevel})
	}
	return texts, nil
}
