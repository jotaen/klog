package terminalformat

type ColourTheme string

const (
	NO_COLOUR = ColourTheme("no_colour")
	DARK      = ColourTheme("dark")
)

func NewStyler(c ColourTheme) Styler {
	baseColouredStyler := Styler{
		props:            StyleProps{},
		colourCodes:      make(map[Colour]string),
		reset:            "\033[0m",
		foregroundPrefix: "\033[38;5;",
		backgroundPrefix: "\033[48;5;",
		colourSuffix:     "m",
		underlined:       "\033[4m",
		bold:             "\033[1m",
	}

	switch c {
	case NO_COLOUR:
		return Styler{
			props:            StyleProps{},
			colourCodes:      make(map[Colour]string),
			reset:            "",
			foregroundPrefix: "",
			backgroundPrefix: "",
			colourSuffix:     "",
			underlined:       "",
			bold:             "",
		}
	case DARK:
		baseColouredStyler.colourCodes = map[Colour]string{
			TEXT:         "015",
			TEXT_INVERSE: "000",
			GREEN:        "120",
			RED:          "167",
			BLUE_DARK:    "117",
			BLUE_LIGHT:   "027",
			SUBDUED:      "249",
			PURPLE:       "213",
			YELLOW:       "221",
		}
		return baseColouredStyler
	}
	panic("Unknown colour theme")
}
