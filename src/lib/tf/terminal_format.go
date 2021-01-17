package tf

type Style struct {
	Color        string
	IsBold       bool
	IsUnderlined bool
}

var reset = "\033[0m"

func (s Style) Format(text string) string {
	return s.seqs() + text + reset
}

func (s Style) FormatAndRestore(text string, previousStyle Style) string {
	return s.Format(text) + previousStyle.seqs()
}

func (s Style) ChangedColor(color string) Style {
	newS := s
	newS.Color = color
	return newS
}

func (s Style) ChangedBold(isBold bool) Style {
	newS := s
	newS.IsBold = isBold
	return newS
}

func (s Style) ChangedUnderlined(isUnderlined bool) Style {
	newS := s
	newS.IsUnderlined = isUnderlined
	return newS
}

func (s Style) seqs() string {
	seqs := reset

	if s.Color != "" {
		seqs = seqs + "\033[38;5;" + s.Color + "m"
	}

	if s.IsUnderlined {
		seqs = seqs + "\033[4m"
	}

	if s.IsBold {
		seqs = seqs + "\033[1m"
	}

	return seqs
}
