package lib

import (
	"errors"
	"fmt"
	"github.com/jotaen/klog/klog/app"
	tf "github.com/jotaen/klog/klog/app/cli/lib/terminalformat"
	"github.com/jotaen/klog/klog/service"
	"strings"
)

var Reflower = tf.NewReflower(60, "\n")

// PrettifyAppError prints app errors including details.
func PrettifyAppError(err app.Error, isDebug bool) error {
	message := "Error: " + err.Error() + "\n"
	message += Reflower.Reflow(err.Details(), nil)
	if isDebug && err.Original() != nil {
		message += "\n\nOriginal Error:\n" + err.Original().Error()
	}
	return errors.New(message)
}

// PrettifyParsingError turns a parsing error into a coloured and well-structured form.
func PrettifyParsingError(err app.ParserErrors, isDebug bool, styler tf.Styler) error {
	message := ""
	INDENT := "    "
	for _, e := range err.All() {
		message += "\n"
		message += fmt.Sprintf(
			styler.Props(tf.StyleProps{Background: tf.RED, Color: tf.RED}).Format("[")+
				styler.Props(tf.StyleProps{Background: tf.RED, Color: tf.TEXT_INVERSE}).Format("SYNTAX ERROR")+
				styler.Props(tf.StyleProps{Background: tf.RED, Color: tf.RED}).Format("]")+
				styler.Props(tf.StyleProps{Color: tf.RED}).Format(" in line %d: "),
			e.LineNumber(),
		) + "\n"
		message += fmt.Sprintf(
			styler.Props(tf.StyleProps{Color: tf.SUBDUED}).Format(INDENT+"%s"),
			// Replace all tabs with one space each, otherwise the carets might
			// not be in line with the text anymore (since we can’t know how wide
			// a tab is).
			strings.Replace(e.LineText(), "\t", " ", -1),
		) + "\n"
		message += fmt.Sprintf(
			styler.Props(tf.StyleProps{Color: tf.RED}).Format(INDENT+"%s%s"),
			strings.Repeat(" ", e.Position()), strings.Repeat("^", e.Length()),
		) + "\n"
		message += fmt.Sprintf(
			styler.Props(tf.StyleProps{Color: tf.YELLOW}).Format("%s"),
			Reflower.Reflow(e.Message(), []string{INDENT}),
		) + "\n"
	}
	return errors.New(message)
}

// PrettifyWarning formats a warning about a record.
func PrettifyWarning(w service.Warning, styler tf.Styler) string {
	return PrettifyGeneralWarning(w.Date().ToString()+": "+w.Warning(), styler)
}

// PrettifyGeneralWarning formats a general warning message.
func PrettifyGeneralWarning(message string, styler tf.Styler) string {
	result := ""
	result += styler.Props(tf.StyleProps{Background: tf.YELLOW, Color: tf.YELLOW}).Format("[")
	result += styler.Props(tf.StyleProps{Background: tf.YELLOW, Color: tf.TEXT_INVERSE}).Format("WARNING")
	result += styler.Props(tf.StyleProps{Background: tf.YELLOW, Color: tf.YELLOW}).Format("]")
	result += " "
	result += styler.Props(tf.StyleProps{Color: tf.YELLOW}).Format(message)
	result += "\n"
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
