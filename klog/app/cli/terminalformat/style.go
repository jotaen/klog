package terminalformat

type StyleProps struct {
	Color        Colour
	Background   Colour
	IsBold       bool
	IsUnderlined bool
}

type Styler struct {
	props            StyleProps
	colourCodes      map[Colour]string
	reset            string
	foregroundPrefix string
	backgroundPrefix string
	colourSuffix     string
	underlined       string
	bold             string
}

type Colour int

const (
	unspecified = iota
	TEXT
	TEXT_INVERSE
	GREEN
	RED
	YELLOW
	BLUE_DARK
	BLUE_LIGHT
	SUBDUED
	PURPLE
)

func (s Styler) Format(text string) string {
	return s.seqs() + text + s.reset
}

func (s Styler) Props(p StyleProps) Styler {
	newS := s
	newS.props = p
	return newS
}

func (s Styler) FormatAndRestore(text string, previousStyle Styler) string {
	return s.Format(text) + previousStyle.seqs()
}

func (s Styler) seqs() string {
	seqs := s.reset

	if s.props.Color != unspecified {
		seqs = seqs + s.foregroundPrefix + s.colourCodes[s.props.Color] + s.colourSuffix
	}

	if s.props.Background != unspecified {
		seqs = seqs + s.backgroundPrefix + s.colourCodes[s.props.Background] + s.colourSuffix
	}

	if s.props.IsUnderlined {
		seqs = seqs + s.underlined
	}

	if s.props.IsBold {
		seqs = seqs + s.bold
	}

	return seqs
}
