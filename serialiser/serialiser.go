package serialiser

import (
	"fmt"
	"klog/workday"
)

func Serialise(workDay workday.WorkDay) string {
	text := ""
	text += fmt.Sprintf("date: %v", workDay.Date().ToString())
	if len(workDay.Summary()) > 0 {
		text += fmt.Sprintf("\nsummary: %v", workDay.Summary())
	}
	text += "\n"
	return text
}
