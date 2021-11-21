package reconciling

import "github.com/jotaen/klog/src/parser/engine"

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
		if len(l.PrecedingWhitespace) > 0 {
			defaultPrefs.indentationStyle = l.PrecedingWhitespace
		}
	}
	return defaultPrefs
}
