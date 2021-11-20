package parsing

// Preferences holds information what kinds of variations were encountered.
// E.g., if the text was indented with tabs or spaces, or if the line endings
// were UNIX or Windows ones.
type Preferences struct {
	LineEnding       string
	IndentationStyle string
}

// DefaultPreferences returns `\n` for the line ending, and 4 spaces for the indentation.
func DefaultPreferences() Preferences {
	return Preferences{
		LineEnding:       "\n",
		IndentationStyle: "    ",
	}
}

// Adapt adjusts the preferences to what is encountered in the Line.
func (p *Preferences) Adapt(l *Line) {
	if len(l.PrecedingWhitespace()) > 0 {
		p.IndentationStyle = l.PrecedingWhitespace()
	}
	if l.originalLineEnding != "" {
		p.LineEnding = l.originalLineEnding
	}
}
