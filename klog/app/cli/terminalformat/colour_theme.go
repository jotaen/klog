package terminalformat

type ColourTheme string
type colourCodes map[Colour]string

const (
	COLOUR_THEME_NO_COLOUR = ColourTheme("no_colour")
	COLOUR_THEME_DARK      = ColourTheme("dark")
	COLOUR_THEME_LIGHT     = ColourTheme("light")
	COLOUR_THEME_BASIC     = ColourTheme("basic")
)

func NewStyler(c ColourTheme) Styler {
	switch c {
	case COLOUR_THEME_NO_COLOUR:
		return Styler{
			props:            StyleProps{},
			colourCodes:      make(colourCodes),
			reset:            "",
			foregroundPrefix: "",
			backgroundPrefix: "",
			colourSuffix:     "",
			underlined:       "",
			bold:             "",
		}
	case COLOUR_THEME_DARK:
		return newStyler256bit(colourCodes{
			TEXT:         "015",
			TEXT_SUBDUED: "249",
			TEXT_INVERSE: "000",
			GREEN:        "120",
			RED:          "167",
			BLUE_DARK:    "117",
			BLUE_LIGHT:   "027",
			PURPLE:       "213",
			YELLOW:       "221",
		})
	case COLOUR_THEME_LIGHT:
		return newStyler256bit(colourCodes{
			TEXT:         "000",
			TEXT_SUBDUED: "237",
			TEXT_INVERSE: "015",
			GREEN:        "028",
			RED:          "124",
			BLUE_DARK:    "025",
			BLUE_LIGHT:   "033",
			PURPLE:       "055",
			YELLOW:       "208",
		})
	case COLOUR_THEME_BASIC:
		return newStyler8bit(colourCodes{
			TEXT:         "", // Disabled
			TEXT_SUBDUED: "", // Disabled
			TEXT_INVERSE: "0",
			GREEN:        "2",
			RED:          "1",
			BLUE_DARK:    "4",
			BLUE_LIGHT:   "6",
			PURPLE:       "5",
			YELLOW:       "3",
		})
	}
	panic("Unknown colour theme")
}

func newStyler256bit(cc colourCodes) Styler {
	return Styler{
		props:            StyleProps{},
		colourCodes:      cc,
		reset:            "\033[0m",
		foregroundPrefix: "\033[38;5;",
		backgroundPrefix: "\033[48;5;",
		colourSuffix:     "m",
		underlined:       "\033[4m",
		bold:             "\033[1m",
	}
}

func newStyler8bit(cc colourCodes) Styler {
	return Styler{
		props:            StyleProps{},
		colourCodes:      cc,
		reset:            "\033[0m",
		foregroundPrefix: "\033[3",
		backgroundPrefix: "\033[4",
		colourSuffix:     "m",
		underlined:       "\033[4m",
		bold:             "\033[1m",
	}
}
