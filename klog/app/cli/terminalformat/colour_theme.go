package terminalformat

type ColourTheme string

const (
	COLOUR_THEME_NO_COLOUR = ColourTheme("no_colour")
	COLOUR_THEME_DARK      = ColourTheme("dark")
	COLOUR_THEME_LIGHT     = ColourTheme("light")
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
	case COLOUR_THEME_NO_COLOUR:
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
	case COLOUR_THEME_DARK:
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
	case COLOUR_THEME_LIGHT:
		baseColouredStyler.colourCodes = map[Colour]string{
			TEXT:         "000",
			TEXT_INVERSE: "015",
			GREEN:        "028",
			RED:          "124",
			BLUE_DARK:    "025",
			BLUE_LIGHT:   "033",
			SUBDUED:      "237",
			PURPLE:       "055",
			YELLOW:       "208",
		}
		return baseColouredStyler
	}
	panic("Unknown colour theme")
}
