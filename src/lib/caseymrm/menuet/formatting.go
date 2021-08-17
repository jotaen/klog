//go:build darwin
// +build darwin

package menuet

// FontWeight represents the weight of the font
type FontWeight float64

const (
	// WeightUltraLight is equivalent to NSFontWeightUltraLight
	WeightUltraLight FontWeight = -0.8
	// WeightThin is equivalent to NSFontWeightThin
	WeightThin = -0.6
	// WeightLight is equivalent to NSFontWeightLight
	WeightLight = -0.4
	// WeightRegular is equivalent to NSFontWeightRegular, and is the default
	WeightRegular = 0
	// WeightMedium is equivalent to NSFontWeightMedium
	WeightMedium = 0.23
	// WeightSemibold is equivalent to NSFontWeightSemibold
	WeightSemibold = 0.3
	// WeightBold is equivalent to NSFontWeightBold
	WeightBold = 0.4
	// WeightHeavy is equivalent to NSFontWeightHeavy
	WeightHeavy = 0.56
	// WeightBlack is equivalent to NSFontWeightBlack
	WeightBlack = 0.62
)
