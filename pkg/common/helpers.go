package common

import "math"

func Similar(a float64, b float64) bool {
	return math.Abs(a-b) < 1e-6
}

func NormalizePlusMinusPi(angle float64) float64 {
	twoPi := math.Pi * 2
	positiveReduced := math.Mod(math.Mod(angle, twoPi)+twoPi, twoPi)

	if positiveReduced > math.Pi {
		positiveReduced -= twoPi
	}

	return positiveReduced
}
