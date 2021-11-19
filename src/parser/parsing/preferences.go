package parsing

// Preferences holds information what kinds of variations were encountered.
// E.g., if the text was indented with tabs or spaces, or if the line endings
// were UNIX or Windows ones.
type Preferences struct {
	LineEnding  string
	Indentation string
}

// DefaultPreferences returns `\n` for the line ending, and 4 spaces for the indentation.
func DefaultPreferences() Preferences {
	return Preferences{
		LineEnding:  "\n",
		Indentation: "    ",
	}
}

// Adapt adjusts the preferences to what is encountered in the Line.
func (p *Preferences) Adapt(l *Line) {
	if l.IndentationLevel() == 1 {
		p.Indentation = l.originalIndentation
	}
	if l.originalLineEnding != "" {
		p.LineEnding = l.originalLineEnding
	}
}
