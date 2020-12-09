package serialiser

import (
	"fmt"
	"klog/workday"
)

func Serialise(workDay workday.WorkDay) string {
	text := ""

	// Date
	text += fmt.Sprintf("date: %v", workDay.Date().ToString())

	// Summary
	if len(workDay.Summary()) > 0 {
		text += fmt.Sprintf("\nsummary: %v", workDay.Summary())
	}

	// Hours
	if len(workDay.Ranges()) > 0 || len(workDay.Times()) > 0 {
		text += "\nhours:"
		for _, rs := range workDay.Ranges() {
			text += fmt.Sprintf("\n- start: %v", rs[0].ToString())
			if rs[1] != nil {
				text += fmt.Sprintf("\n  end: %v", rs[1].ToString())
			}
		}
		for _, t := range workDay.Times() {
			text += fmt.Sprintf("\n- time: %v", t.ToString())
		}
	}

	// Final newline
	text += "\n"
	return text
}
