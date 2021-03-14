package parsing

type Preferences struct {
	LineEnding  string
	Indentation string
}

func DefaultPreferences() Preferences {
	return Preferences{
		LineEnding:  "\n",
		Indentation: "    ",
	}
}

func (p *Preferences) Adapt(l *Line) {
	if l.IndentationLevel() == 1 {
		p.Indentation = l.originalIndentation
	}
	if l.originalLineEnding != "" {
		p.LineEnding = l.originalLineEnding
	}
}
