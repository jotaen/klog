package reconciling

import (
	"github.com/jotaen/klog/src/parser/engine"
	"strings"
)

type stylePreferences struct {
	indentationStyle string
	lineEndingStyle  string
}

func stylePreferencesOrDefault(b engine.Block) stylePreferences {
	defaultPrefs := stylePreferences{
		indentationStyle: "    ",
		lineEndingStyle:  "\n",
	}
	if b == nil {
		return defaultPrefs
	}
	for _, l := range b.SignificantLines() {
		if len(l.LineEnding) > 0 {
			defaultPrefs.lineEndingStyle = l.LineEnding
		}
		precedingWhitespace := precedingWhitespace(l.Text)
		if len(precedingWhitespace) > 0 {
			defaultPrefs.indentationStyle = precedingWhitespace
		}
	}
	return defaultPrefs
}

func precedingWhitespace(line string) string {
	netText := strings.TrimLeftFunc(line, engine.IsSpaceOrTab)
	return line[:len(line)-len(netText)]
}
