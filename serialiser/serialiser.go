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
	hasHours := len(workDay.Ranges()) > 0 || len(workDay.Times()) > 0
	if hasHours {
		text += "\nhours:"
		for _, r := range workDay.Ranges() {
			text += fmt.Sprintf("\n- start: %v", r.Start().ToString())
			if !r.IsOpen() {
				text += fmt.Sprintf("\n  end: %v", r.End().ToString())
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
