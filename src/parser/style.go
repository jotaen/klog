package parser

// Style describes the general styling and formatting preferences of a record.
type Style struct {
	LineEnding       string
	Indentation      string
	UsesDashesInDate bool // Example: 2000-01-01
	Uses24HourClock  bool // Example: 8:00
	UsesSpaceInRange bool // Example: 8:00 - 9:00
}

// DefaultStyle returns the canonical style preferences as recommended
// by the file format specification.
func DefaultStyle() Style {
	return Style{
		LineEnding:       "\n",
		Indentation:      "    ",
		UsesDashesInDate: true,
		Uses24HourClock:  true,
		UsesSpaceInRange: true,
	}
}
