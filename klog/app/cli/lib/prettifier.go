package lib

import (
	"errors"
	"fmt"
	"github.com/jotaen/klog/klog/app"
	"github.com/jotaen/klog/klog/app/cli/lib/terminalformat"
	"github.com/jotaen/klog/klog/service"
	"strings"
)

// PrettifyError turns an error into a coloured and well-structured form.
func PrettifyError(err error, isDebug bool) error {
	reflower := terminalformat.NewReflower(60, "\n")
	switch e := err.(type) {
	case app.ParserErrors:
		message := ""
		INDENT := "    "
		for _, e := range e.All() {
			message += fmt.Sprintf(
				terminalformat.Style{Background: "160", Color: "015"}.Format(" ERROR in line %d: "),
				e.LineNumber(),
			) + "\n"
			message += fmt.Sprintf(
				terminalformat.Style{Color: "247"}.Format(INDENT+"%s"),
				// Replace all tabs with one space each, otherwise the carets might
				// not be in line with the text anymore (since we can’t know how wide
				// a tab is).
				strings.Replace(e.LineText(), "\t", " ", -1),
			) + "\n"
			message += fmt.Sprintf(
				terminalformat.Style{Color: "160"}.Format(INDENT+"%s%s"),
				strings.Repeat(" ", e.Position()), strings.Repeat("^", e.Length()),
			) + "\n"
			message += fmt.Sprintf(
				terminalformat.Style{Color: "227"}.Format("%s"),
				reflower.Reflow(e.Message(), INDENT),
			) + "\n\n"
		}
		return errors.New(message)
	case app.Error:
		message := "Error: " + e.Error() + "\n"
		message += reflower.Reflow(e.Details(), "")
		if isDebug && e.Original() != nil {
			message += "\n\nOriginal Error:\n" + e.Original().Error()
		}
		return errors.New(message)
	}
	return errors.New("Error: " + err.Error())
}

// PrettifyWarnings turns an error into a coloured and well-structured form.
func PrettifyWarnings(ws []service.Warning) string {
	result := ""
	for _, w := range ws {
		result += terminalformat.Style{Background: "227", Color: "000"}.Format(" WARNING ")
		result += " "
		result += terminalformat.Style{Color: "227"}.Format(w.Date().ToString() + ": " + w.Warning())
		result += "\n"
	}
	return result
}

// PrettyMonth returns the full english name of a month.
func PrettyMonth(m int) string {
	switch m {
	case 1:
		return "January"
	case 2:
		return "February"
	case 3:
		return "March"
	case 4:
		return "April"
	case 5:
		return "May"
	case 6:
		return "June"
	case 7:
		return "July"
	case 8:
		return "August"
	case 9:
		return "September"
	case 10:
		return "October"
	case 11:
		return "November"
	case 12:
		return "December"
	}
	panic("Illegal month") // this can/should never happen
}

// PrettyDay returns the full english name of a weekday.
func PrettyDay(d int) string {
	switch d {
	case 1:
		return "Monday"
	case 2:
		return "Tuesday"
	case 3:
		return "Wednesday"
	case 4:
		return "Thursday"
	case 5:
		return "Friday"
	case 6:
		return "Saturday"
	case 7:
		return "Sunday"
	}
	panic("Illegal weekday") // this can/should never happen
}
